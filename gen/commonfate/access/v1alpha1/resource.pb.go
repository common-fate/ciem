// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: commonfate/access/v1alpha1/resource.proto

package accessv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Resource:
	//
	//	*Resource_AwsAccount
	//	*Resource_GcpProject
	Resource isResource_Resource `protobuf_oneof:"resource"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_commonfate_access_v1alpha1_resource_proto_rawDescGZIP(), []int{0}
}

func (m *Resource) GetResource() isResource_Resource {
	if m != nil {
		return m.Resource
	}
	return nil
}

func (x *Resource) GetAwsAccount() *AWSAccount {
	if x, ok := x.GetResource().(*Resource_AwsAccount); ok {
		return x.AwsAccount
	}
	return nil
}

func (x *Resource) GetGcpProject() *GCPProject {
	if x, ok := x.GetResource().(*Resource_GcpProject); ok {
		return x.GcpProject
	}
	return nil
}

type isResource_Resource interface {
	isResource_Resource()
}

type Resource_AwsAccount struct {
	AwsAccount *AWSAccount `protobuf:"bytes,1,opt,name=aws_account,json=awsAccount,proto3,oneof"`
}

type Resource_GcpProject struct {
	GcpProject *GCPProject `protobuf:"bytes,2,opt,name=gcp_project,json=gcpProject,proto3,oneof"`
}

func (*Resource_AwsAccount) isResource_Resource() {}

func (*Resource_GcpProject) isResource_Resource() {}

type GCPProject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project string `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Role    string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *GCPProject) Reset() {
	*x = GCPProject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GCPProject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GCPProject) ProtoMessage() {}

func (x *GCPProject) ProtoReflect() protoreflect.Message {
	mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GCPProject.ProtoReflect.Descriptor instead.
func (*GCPProject) Descriptor() ([]byte, []int) {
	return file_commonfate_access_v1alpha1_resource_proto_rawDescGZIP(), []int{1}
}

func (x *GCPProject) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *GCPProject) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type AWSAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Role      string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *AWSAccount) Reset() {
	*x = AWSAccount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AWSAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AWSAccount) ProtoMessage() {}

func (x *AWSAccount) ProtoReflect() protoreflect.Message {
	mi := &file_commonfate_access_v1alpha1_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AWSAccount.ProtoReflect.Descriptor instead.
func (*AWSAccount) Descriptor() ([]byte, []int) {
	return file_commonfate_access_v1alpha1_resource_proto_rawDescGZIP(), []int{2}
}

func (x *AWSAccount) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *AWSAccount) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

var File_commonfate_access_v1alpha1_resource_proto protoreflect.FileDescriptor

var file_commonfate_access_v1alpha1_resource_proto_rawDesc = []byte{
	0x0a, 0x29, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2f, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x20, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0xb8, 0x01,
	0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0b, 0x61, 0x77,
	0x73, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x41, 0x57, 0x53, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x48, 0x00, 0x52,
	0x0a, 0x61, 0x77, 0x73, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x4f, 0x0a, 0x0b, 0x67,
	0x63, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x47, 0x43, 0x50, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x48, 0x00,
	0x52, 0x0a, 0x67, 0x63, 0x70, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x0a, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x3a, 0x0a, 0x0a, 0x47, 0x43, 0x50, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x72, 0x6f, 0x6c, 0x65, 0x22, 0x3f, 0x0a, 0x0a, 0x41, 0x57, 0x53, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x42, 0xa3, 0x02, 0x0a, 0x24, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0d,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2d, 0x66, 0x61, 0x74, 0x65, 0x2f, 0x63, 0x69, 0x65, 0x6d, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2f, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x43, 0x41,
	0xaa, 0x02, 0x20, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65, 0x2e, 0x43, 0x6c,
	0x6f, 0x75, 0x64, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0xca, 0x02, 0x20, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61, 0x74, 0x65,
	0x5c, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x5c, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x2c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66,
	0x61, 0x74, 0x65, 0x5c, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x5c, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x23, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x66, 0x61,
	0x74, 0x65, 0x3a, 0x3a, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x3a, 0x3a, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_commonfate_access_v1alpha1_resource_proto_rawDescOnce sync.Once
	file_commonfate_access_v1alpha1_resource_proto_rawDescData = file_commonfate_access_v1alpha1_resource_proto_rawDesc
)

func file_commonfate_access_v1alpha1_resource_proto_rawDescGZIP() []byte {
	file_commonfate_access_v1alpha1_resource_proto_rawDescOnce.Do(func() {
		file_commonfate_access_v1alpha1_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_commonfate_access_v1alpha1_resource_proto_rawDescData)
	})
	return file_commonfate_access_v1alpha1_resource_proto_rawDescData
}

var file_commonfate_access_v1alpha1_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_commonfate_access_v1alpha1_resource_proto_goTypes = []interface{}{
	(*Resource)(nil),   // 0: commonfate.cloud.access.v1alpha1.Resource
	(*GCPProject)(nil), // 1: commonfate.cloud.access.v1alpha1.GCPProject
	(*AWSAccount)(nil), // 2: commonfate.cloud.access.v1alpha1.AWSAccount
}
var file_commonfate_access_v1alpha1_resource_proto_depIdxs = []int32{
	2, // 0: commonfate.cloud.access.v1alpha1.Resource.aws_account:type_name -> commonfate.cloud.access.v1alpha1.AWSAccount
	1, // 1: commonfate.cloud.access.v1alpha1.Resource.gcp_project:type_name -> commonfate.cloud.access.v1alpha1.GCPProject
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_commonfate_access_v1alpha1_resource_proto_init() }
func file_commonfate_access_v1alpha1_resource_proto_init() {
	if File_commonfate_access_v1alpha1_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_commonfate_access_v1alpha1_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commonfate_access_v1alpha1_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GCPProject); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commonfate_access_v1alpha1_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AWSAccount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_commonfate_access_v1alpha1_resource_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Resource_AwsAccount)(nil),
		(*Resource_GcpProject)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commonfate_access_v1alpha1_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commonfate_access_v1alpha1_resource_proto_goTypes,
		DependencyIndexes: file_commonfate_access_v1alpha1_resource_proto_depIdxs,
		MessageInfos:      file_commonfate_access_v1alpha1_resource_proto_msgTypes,
	}.Build()
	File_commonfate_access_v1alpha1_resource_proto = out.File
	file_commonfate_access_v1alpha1_resource_proto_rawDesc = nil
	file_commonfate_access_v1alpha1_resource_proto_goTypes = nil
	file_commonfate_access_v1alpha1_resource_proto_depIdxs = nil
}