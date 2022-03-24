// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: arrebato/acl/v1/acl.proto

// Package arrebato.acl.v1 provides messages that describe access-control lists (ACLs) for clients wishing to consume
// and produce messages on individual topics.

package acl

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

// The Permission enumeration represents individual permissions a client may have on
// a topic.
type Permission int32

const (
	// The default permission value, using this should always result in an error.
	Permission_PERMISSION_UNSPECIFIED Permission = 0
	// The client may produce messages on the topic.
	Permission_PERMISSION_PRODUCE Permission = 1
	// The client may consume messages from the topic.
	Permission_PERMISSION_CONSUME Permission = 2
)

// Enum value maps for Permission.
var (
	Permission_name = map[int32]string{
		0: "PERMISSION_UNSPECIFIED",
		1: "PERMISSION_PRODUCE",
		2: "PERMISSION_CONSUME",
	}
	Permission_value = map[string]int32{
		"PERMISSION_UNSPECIFIED": 0,
		"PERMISSION_PRODUCE":     1,
		"PERMISSION_CONSUME":     2,
	}
)

func (x Permission) Enum() *Permission {
	p := new(Permission)
	*p = x
	return p
}

func (x Permission) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Permission) Descriptor() protoreflect.EnumDescriptor {
	return file_arrebato_acl_v1_acl_proto_enumTypes[0].Descriptor()
}

func (Permission) Type() protoreflect.EnumType {
	return &file_arrebato_acl_v1_acl_proto_enumTypes[0]
}

func (x Permission) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Permission.Descriptor instead.
func (Permission) EnumDescriptor() ([]byte, []int) {
	return file_arrebato_acl_v1_acl_proto_rawDescGZIP(), []int{0}
}

// The ACL message describes the entire state of the access-control list.
type ACL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The permissions for each client/topic pair.
	Entries []*Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *ACL) Reset() {
	*x = ACL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_acl_v1_acl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ACL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ACL) ProtoMessage() {}

func (x *ACL) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_acl_v1_acl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ACL.ProtoReflect.Descriptor instead.
func (*ACL) Descriptor() ([]byte, []int) {
	return file_arrebato_acl_v1_acl_proto_rawDescGZIP(), []int{0}
}

func (x *ACL) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

// The Entry message describes a single set of permissions applied to a client
// for a topic.
type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The topic to set permissions for.
	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	// The client identifier the entry refers to. In an insecure environment, this can be an arbitrary string that the
	// client will use to identify itself in the request metadata. When using mutual TLS, this will be a SPIFFE ID that
	// the client will include in its TLS certificate.
	Client string `protobuf:"bytes,2,opt,name=client,proto3" json:"client,omitempty"`
	// Permissions to apply to the client.
	Permissions []Permission `protobuf:"varint,3,rep,packed,name=permissions,proto3,enum=arrebato.acl.v1.Permission" json:"permissions,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_acl_v1_acl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_acl_v1_acl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_arrebato_acl_v1_acl_proto_rawDescGZIP(), []int{1}
}

func (x *Entry) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Entry) GetClient() string {
	if x != nil {
		return x.Client
	}
	return ""
}

func (x *Entry) GetPermissions() []Permission {
	if x != nil {
		return x.Permissions
	}
	return nil
}

var File_arrebato_acl_v1_acl_proto protoreflect.FileDescriptor

var file_arrebato_acl_v1_acl_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x6c, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x63, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x72, 0x72,
	0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x61, 0x63, 0x6c, 0x2e, 0x76, 0x31, 0x22, 0x37, 0x0a, 0x03,
	0x41, 0x43, 0x4c, 0x12, 0x30, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e,
	0x61, 0x63, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x74, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x3d, 0x0a, 0x0b,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0e, 0x32, 0x1b, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x61, 0x63, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2a, 0x58, 0x0a, 0x0a, 0x50,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x45, 0x52,
	0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x44, 0x55, 0x43, 0x45, 0x10, 0x01, 0x12, 0x16, 0x0a,
	0x12, 0x50, 0x45, 0x52, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x4f, 0x4e, 0x53,
	0x55, 0x4d, 0x45, 0x10, 0x02, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76, 0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e, 0x64, 0x2f, 0x61,
	0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f,
	0x61, 0x63, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_arrebato_acl_v1_acl_proto_rawDescOnce sync.Once
	file_arrebato_acl_v1_acl_proto_rawDescData = file_arrebato_acl_v1_acl_proto_rawDesc
)

func file_arrebato_acl_v1_acl_proto_rawDescGZIP() []byte {
	file_arrebato_acl_v1_acl_proto_rawDescOnce.Do(func() {
		file_arrebato_acl_v1_acl_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_acl_v1_acl_proto_rawDescData)
	})
	return file_arrebato_acl_v1_acl_proto_rawDescData
}

var file_arrebato_acl_v1_acl_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_arrebato_acl_v1_acl_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_arrebato_acl_v1_acl_proto_goTypes = []interface{}{
	(Permission)(0), // 0: arrebato.acl.v1.Permission
	(*ACL)(nil),     // 1: arrebato.acl.v1.ACL
	(*Entry)(nil),   // 2: arrebato.acl.v1.Entry
}
var file_arrebato_acl_v1_acl_proto_depIdxs = []int32{
	2, // 0: arrebato.acl.v1.ACL.entries:type_name -> arrebato.acl.v1.Entry
	0, // 1: arrebato.acl.v1.Entry.permissions:type_name -> arrebato.acl.v1.Permission
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_arrebato_acl_v1_acl_proto_init() }
func file_arrebato_acl_v1_acl_proto_init() {
	if File_arrebato_acl_v1_acl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_acl_v1_acl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ACL); i {
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
		file_arrebato_acl_v1_acl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_arrebato_acl_v1_acl_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arrebato_acl_v1_acl_proto_goTypes,
		DependencyIndexes: file_arrebato_acl_v1_acl_proto_depIdxs,
		EnumInfos:         file_arrebato_acl_v1_acl_proto_enumTypes,
		MessageInfos:      file_arrebato_acl_v1_acl_proto_msgTypes,
	}.Build()
	File_arrebato_acl_v1_acl_proto = out.File
	file_arrebato_acl_v1_acl_proto_rawDesc = nil
	file_arrebato_acl_v1_acl_proto_goTypes = nil
	file_arrebato_acl_v1_acl_proto_depIdxs = nil
}
