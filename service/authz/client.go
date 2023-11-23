package authz

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bufbuild/connect-go"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1/authzv1alpha1connect"
	"github.com/fatih/structtag"
)

type Client struct {
	// Raw returns the underlying GRPC client.
	// It can be used to call methods that we don't have wrappers for yet.
	raw authzv1alpha1connect.AuthzServiceClient
}

func NewClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) Client {
	// client ALWAYS uses GRPC as Rust does not support Buf Connect.
	opts = append(opts, connect.WithGRPC())
	return Client{raw: authzv1alpha1connect.NewAuthzServiceClient(httpClient, baseURL, opts...)}
}

// Entities are objects that can be stored in the authz database.
type Entity interface {
	EntityType() string
}

func transformToEntity(e Entity) (*authzv1alpha1.Entity, error) {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct", v.Type())
	}

	entity := authzv1alpha1.Entity{
		Attributes: []*authzv1alpha1.Attribute{},
		Parents:    []*authzv1alpha1.UID{},
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		t, err := parseTag(string(f.Tag))
		if err != nil {
			return nil, err
		}

		if t.Name == "id" {
			switch val := v.Field(i); val.Kind() {
			case reflect.String:
				entity.Uid = &authzv1alpha1.UID{
					Type: e.EntityType(),
					Id:   val.String(),
				}

			default:
				return nil, errors.New("unsupported ID field type, only string IDs are currently supported")
			}
			continue
		}

		// try and parse as a parent
		if t.ParentType != "" {
			switch val := v.Field(i); val.Kind() {
			case reflect.String:
				entity.Parents = append(entity.Parents, &authzv1alpha1.UID{
					Type: t.ParentType,
					Id:   val.String(),
				})

			case reflect.Slice:
				slice, ok := val.Interface().([]string)
				if !ok {
					return nil, errors.New("invalid slice: unsupported parent field type, only strings and string slice IDs are currently supported")
				}

				for _, s := range slice {
					entity.Parents = append(entity.Parents, &authzv1alpha1.UID{
						Type: t.ParentType,
						Id:   s,
					})
				}

			default:
				return nil, errors.New("unsupported parent field type, only strings and string slice IDs are currently supported")
			}
			continue
		}

		// try and pass as a generic parent
		if t.HasParent {
			val := v.Field(i)
			slice, ok := val.Interface().(*authzv1alpha1.UID)
			if ok {
				entity.Parents = append(entity.Parents, slice)
				continue
			}

			switch val.Kind() {
			case reflect.Slice:
				slice, ok := val.Interface().([]*authzv1alpha1.UID)
				if !ok {
					return nil, errors.New("invalid slice: unsupported parent field type, []*authzv1alpha1.UID slices are currently supported")
				}

				entity.Parents = append(entity.Parents, slice...)

			default:
				return nil, errors.New("unsupported parent field type, only []*authzv1alpha1.UID slices are currently supported. Otherwise, specify a particular parent entity type with parent=<type>")
			}
			continue
		}

		// try and parse as an attribute
		var attr *authzv1alpha1.Value

		switch val := v.Field(i); val.Kind() {
		case reflect.String:
			attr = &authzv1alpha1.Value{
				Value: &authzv1alpha1.Value_Str{
					Str: val.String(),
				},
			}
		case reflect.Bool:
			attr = &authzv1alpha1.Value{
				Value: &authzv1alpha1.Value_Bool{
					Bool: val.Bool(),
				},
			}

		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			attr = &authzv1alpha1.Value{
				Value: &authzv1alpha1.Value_Long{
					Long: val.Int(),
				},
			}

		default:
			return nil, fmt.Errorf("unsupported attribute field type: %s (tag %s)", f.Type.Name(), t.Name)
		}

		entity.Attributes = append(entity.Attributes, &authzv1alpha1.Attribute{
			Key:   t.Name,
			Value: attr,
		})
	}

	return &entity, nil
}

type Tag struct {
	// Name of the tag, e.g. in `authz:"id"` it is "id"
	Name string

	// the "parent" key e.g. in  `authz:"groups,parent=Group"` it is "Group"
	ParentType string

	// Parent is true if the struct tag is like `authz:"resources,parent"`
	HasParent bool
}

func parseTag(input string) (Tag, error) {
	tags, err := structtag.Parse(input)
	if err != nil {
		return Tag{}, err
	}
	authzTag, err := tags.Get("authz")
	if err != nil {
		return Tag{}, err
	}

	t := Tag{
		Name:       authzTag.Name,
		ParentType: extractParentType(authzTag),
		HasParent:  strings.Contains(input, "parent"),
	}

	return t, nil
}

// extractParentType extracts the parent option from the struct tag.
// e.g. if the tag is `authz:"groups,parent=Group"` the field will be "Group".
func extractParentType(t *structtag.Tag) string {
	for _, opt := range t.Options {
		// split "type=User" into ["type", "User"]
		splits := strings.Split(opt, "=")
		if len(splits) < 2 {
			continue
		}
		if splits[0] == "parent" {
			return splits[1]
		}
	}

	return ""
}

func UnmarshalEntity(e *authzv1alpha1.Entity, out Entity) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("output must be a pointer to a struct")
	}

	// Check EntityType matches the UID Type
	if e.Uid == nil || e.Uid.Type != out.EntityType() {
		return fmt.Errorf("entity type mismatch: expected %s, got %s", out.EntityType(), e.Uid.Type)
	}

	v = v.Elem()

	// Handle UID
	if e.Uid != nil && e.Uid.Id != "" {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			tag, _ := parseTag(string(field.Tag))

			if tag.Name == "id" {
				v.Field(i).SetString(e.Uid.Id)
				break
			}
		}
	}

	// Handle Attributes
	for _, attr := range e.Attributes {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			tag, _ := parseTag(string(field.Tag))

			if tag.Name == attr.Key {
				fieldValue := v.Field(i)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(attr.Value.GetStr())
				case reflect.Bool:
					fieldValue.SetBool(attr.Value.GetBool())
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
					fieldValue.SetInt(attr.Value.GetLong())
					// Add other cases as necessary
				}
				break
			}
		}
	}

	// Handle Parents
	// Assuming Parents are mapped to fields marked with `parent` tag
	for _, parent := range e.Parents {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			tag, _ := parseTag(string(field.Tag))

			if tag.ParentType != "" && parent.Type == tag.ParentType {
				fieldValue := v.Field(i)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(parent.Id)
				case reflect.Slice:
					// Check if the slice is of type string
					if fieldValue.Type().Elem().Kind() == reflect.String {
						// Append the new string to the slice
						updatedSlice := reflect.Append(fieldValue, reflect.ValueOf(parent.Id))
						fieldValue.Set(updatedSlice)
					}
				}
				break
			}
		}
	}

	return nil
}
