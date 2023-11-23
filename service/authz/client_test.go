package authz

import (
	"reflect"
	"testing"

	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/stretchr/testify/assert"
)

type testUser struct {
	ID     string   `authz:"id"`
	Name   string   `authz:"name"`
	Groups []string `authz:"groups,parent=Group"`
}

func (testUser) EntityType() string { return "User" }

type testAccount struct {
	ID      string `authz:"id"`
	Name    string `authz:"name"`
	OrgUnit string `authz:"org_unit,parent=OrgUnit"`
}

func (testAccount) EntityType() string { return "Account" }

type testVault struct {
	ID      string `authz:"id"`
	LongVal int    `authz:"long_val"`
	BoolVal bool   `authz:"bool_val"`
}

func (testVault) EntityType() string { return "Vault" }

type testAnyParents struct {
	ID        string               `authz:"id"`
	Resources []*authzv1alpha1.UID `authz:"resources,parent"`
}

func (testAnyParents) EntityType() string { return "AnyParent" }

type testAccessRequest struct {
	ID       string             `authz:"id"`
	Resource *authzv1alpha1.UID `authz:"resource,parent"`
}

func (testAccessRequest) EntityType() string { return "AccessRequest" }

func Test_transformToEntity(t *testing.T) {
	type args struct {
		e Entity
	}
	tests := []struct {
		name    string
		args    args
		want    *authzv1alpha1.Entity
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				e: testUser{
					ID:     "test",
					Name:   "testing",
					Groups: []string{"devs"},
				},
			},
			want: &authzv1alpha1.Entity{
				Uid: &authzv1alpha1.UID{
					Type: "User",
					Id:   "test",
				},
				Attributes: []*authzv1alpha1.Attribute{
					{
						Key: "name",
						Value: &authzv1alpha1.Value{
							Value: &authzv1alpha1.Value_Str{
								Str: "testing",
							},
						},
					},
				},
				Parents: []*authzv1alpha1.UID{
					{
						Type: "Group",
						Id:   "devs",
					},
				},
			},
		},
		{
			name: "many to one",
			args: args{
				e: testAccount{
					ID:      "test",
					Name:    "testing",
					OrgUnit: "prod",
				},
			},
			want: &authzv1alpha1.Entity{
				Uid: &authzv1alpha1.UID{
					Type: "Account",
					Id:   "test",
				},
				Attributes: []*authzv1alpha1.Attribute{
					{
						Key: "name",
						Value: &authzv1alpha1.Value{
							Value: &authzv1alpha1.Value_Str{
								Str: "testing",
							},
						},
					},
				},
				Parents: []*authzv1alpha1.UID{
					{
						Type: "OrgUnit",
						Id:   "prod",
					},
				},
			},
		},
		{
			name: "attribute parsing",
			args: args{
				e: testVault{
					ID:      "test",
					LongVal: 1,
					BoolVal: true,
				},
			},
			want: &authzv1alpha1.Entity{
				Uid: &authzv1alpha1.UID{
					Type: "Vault",
					Id:   "test",
				},
				Attributes: []*authzv1alpha1.Attribute{
					{
						Key: "long_val",
						Value: &authzv1alpha1.Value{
							Value: &authzv1alpha1.Value_Long{
								Long: 1,
							},
						},
					},
					{
						Key: "bool_val",
						Value: &authzv1alpha1.Value{
							Value: &authzv1alpha1.Value_Bool{
								Bool: true,
							},
						},
					},
				},
				Parents: []*authzv1alpha1.UID{},
			},
		},
		{
			name: "generic parents",
			args: args{
				e: testAnyParents{
					ID: "test",
					Resources: []*authzv1alpha1.UID{
						{
							Type: "Something",
							Id:   "test",
						},
						{
							Type: "Other",
							Id:   "else",
						},
					},
				},
			},
			want: &authzv1alpha1.Entity{
				Uid: &authzv1alpha1.UID{
					Type: "AnyParent",
					Id:   "test",
				},
				Attributes: []*authzv1alpha1.Attribute{},
				Parents: []*authzv1alpha1.UID{
					{
						Type: "Something",
						Id:   "test",
					},
					{
						Type: "Other",
						Id:   "else",
					},
				},
			},
		},
		{
			name: "access request",
			args: args{
				e: testAccessRequest{
					ID: "test",
					Resource: &authzv1alpha1.UID{
						Type: "Something",
						Id:   "test",
					},
				},
			},
			want: &authzv1alpha1.Entity{
				Uid: &authzv1alpha1.UID{
					Type: "AccessRequest",
					Id:   "test",
				},
				Attributes: []*authzv1alpha1.Attribute{},
				Parents: []*authzv1alpha1.UID{
					{
						Type: "Something",
						Id:   "test",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformToEntity(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseTag(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    Tag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				input: `authz:"id"`,
			},
			want: Tag{
				Name: "id",
			},
		},
		{
			name: "parent",
			args: args{
				input: `authz:"group,parent=Group"`,
			},
			want: Tag{
				Name:       "group",
				ParentType: "Group",
				HasParent:  true,
			},
		},
		{
			name: "with other tags",
			args: args{
				input: `authz:"group,parent=Group" json:"something"`,
			},
			want: Tag{
				Name:       "group",
				ParentType: "Group",
				HasParent:  true,
			},
		},
		{
			name: "with generic parent",
			args: args{
				input: `authz:"group,parent"`,
			},
			want: Tag{
				Name:      "group",
				HasParent: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTag(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
