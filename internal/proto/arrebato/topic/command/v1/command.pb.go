// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: arrebato/topic/command/v1/command.proto

// Package arrebato.topic.command.v1 provides commands written to the raft log regarding topic management. These messages
// are used internally by nodes in the cluster to portray changes from the leader node to all followers.

package topiccmd

import (
	v1 "github.com/davidsbond/arrebato/internal/proto/arrebato/topic/v1"
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

// The CreateTopic message is a command indicating that a new topic has been created.
type CreateTopic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The new topic.
	Topic *v1.Topic `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *CreateTopic) Reset() {
	*x = CreateTopic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_topic_command_v1_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTopic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTopic) ProtoMessage() {}

func (x *CreateTopic) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_topic_command_v1_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTopic.ProtoReflect.Descriptor instead.
func (*CreateTopic) Descriptor() ([]byte, []int) {
	return file_arrebato_topic_command_v1_command_proto_rawDescGZIP(), []int{0}
}

func (x *CreateTopic) GetTopic() *v1.Topic {
	if x != nil {
		return x.Topic
	}
	return nil
}

// The DeleteTopic message is a command indicating that a topic should be deleted.
type DeleteTopic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the topic to delete.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteTopic) Reset() {
	*x = DeleteTopic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_topic_command_v1_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTopic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTopic) ProtoMessage() {}

func (x *DeleteTopic) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_topic_command_v1_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTopic.ProtoReflect.Descriptor instead.
func (*DeleteTopic) Descriptor() ([]byte, []int) {
	return file_arrebato_topic_command_v1_command_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteTopic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_arrebato_topic_command_v1_command_proto protoreflect.FileDescriptor

var file_arrebato_topic_command_v1_command_proto_rawDesc = []byte{
	0x0a, 0x27, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x61, 0x72, 0x72, 0x65, 0x62,
	0x61, 0x74, 0x6f, 0x2e, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x76, 0x31, 0x1a, 0x1d, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x70,
	0x69, 0x63, 0x12, 0x2e, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x22, 0x21, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x70, 0x69,
	0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x52, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76, 0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e, 0x64, 0x2f, 0x61,
	0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x76, 0x31,
	0x3b, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x63, 0x6d, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_arrebato_topic_command_v1_command_proto_rawDescOnce sync.Once
	file_arrebato_topic_command_v1_command_proto_rawDescData = file_arrebato_topic_command_v1_command_proto_rawDesc
)

func file_arrebato_topic_command_v1_command_proto_rawDescGZIP() []byte {
	file_arrebato_topic_command_v1_command_proto_rawDescOnce.Do(func() {
		file_arrebato_topic_command_v1_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_topic_command_v1_command_proto_rawDescData)
	})
	return file_arrebato_topic_command_v1_command_proto_rawDescData
}

var file_arrebato_topic_command_v1_command_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_arrebato_topic_command_v1_command_proto_goTypes = []interface{}{
	(*CreateTopic)(nil), // 0: arrebato.topic.command.v1.CreateTopic
	(*DeleteTopic)(nil), // 1: arrebato.topic.command.v1.DeleteTopic
	(*v1.Topic)(nil),    // 2: arrebato.topic.v1.Topic
}
var file_arrebato_topic_command_v1_command_proto_depIdxs = []int32{
	2, // 0: arrebato.topic.command.v1.CreateTopic.topic:type_name -> arrebato.topic.v1.Topic
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_arrebato_topic_command_v1_command_proto_init() }
func file_arrebato_topic_command_v1_command_proto_init() {
	if File_arrebato_topic_command_v1_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_topic_command_v1_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTopic); i {
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
		file_arrebato_topic_command_v1_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTopic); i {
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
			RawDescriptor: file_arrebato_topic_command_v1_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arrebato_topic_command_v1_command_proto_goTypes,
		DependencyIndexes: file_arrebato_topic_command_v1_command_proto_depIdxs,
		MessageInfos:      file_arrebato_topic_command_v1_command_proto_msgTypes,
	}.Build()
	File_arrebato_topic_command_v1_command_proto = out.File
	file_arrebato_topic_command_v1_command_proto_rawDesc = nil
	file_arrebato_topic_command_v1_command_proto_goTypes = nil
	file_arrebato_topic_command_v1_command_proto_depIdxs = nil
}
