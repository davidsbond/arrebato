// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: arrebato/message/command/v1/command.proto

// Package arrebato.message.command.v1 provides commands written to the raft log regarding produced/consumed messages.
// These messages are used internally by nodes in the cluster to portray changes from the leader node to all followers.

package messagecmd

import (
	v1 "github.com/davidsbond/arrebato/internal/proto/arrebato/message/v1"
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

// The CreateMessage message is a command indicating that a new message has been produced.
type CreateMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The message being produced.
	Message *v1.Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateMessage) Reset() {
	*x = CreateMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_message_command_v1_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMessage) ProtoMessage() {}

func (x *CreateMessage) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_message_command_v1_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMessage.ProtoReflect.Descriptor instead.
func (*CreateMessage) Descriptor() ([]byte, []int) {
	return file_arrebato_message_command_v1_command_proto_rawDescGZIP(), []int{0}
}

func (x *CreateMessage) GetMessage() *v1.Message {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_arrebato_message_command_v1_command_proto protoreflect.FileDescriptor

var file_arrebato_message_command_v1_command_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x72, 0x72,
	0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61,
	0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x0d, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x42, 0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76, 0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e, 0x64, 0x2f, 0x61, 0x72,
	0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x76,
	0x31, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x63, 0x6d, 0x64, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arrebato_message_command_v1_command_proto_rawDescOnce sync.Once
	file_arrebato_message_command_v1_command_proto_rawDescData = file_arrebato_message_command_v1_command_proto_rawDesc
)

func file_arrebato_message_command_v1_command_proto_rawDescGZIP() []byte {
	file_arrebato_message_command_v1_command_proto_rawDescOnce.Do(func() {
		file_arrebato_message_command_v1_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_message_command_v1_command_proto_rawDescData)
	})
	return file_arrebato_message_command_v1_command_proto_rawDescData
}

var file_arrebato_message_command_v1_command_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_arrebato_message_command_v1_command_proto_goTypes = []interface{}{
	(*CreateMessage)(nil), // 0: arrebato.message.command.v1.CreateMessage
	(*v1.Message)(nil),    // 1: arrebato.message.v1.Message
}
var file_arrebato_message_command_v1_command_proto_depIdxs = []int32{
	1, // 0: arrebato.message.command.v1.CreateMessage.message:type_name -> arrebato.message.v1.Message
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_arrebato_message_command_v1_command_proto_init() }
func file_arrebato_message_command_v1_command_proto_init() {
	if File_arrebato_message_command_v1_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_message_command_v1_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMessage); i {
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
			RawDescriptor: file_arrebato_message_command_v1_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arrebato_message_command_v1_command_proto_goTypes,
		DependencyIndexes: file_arrebato_message_command_v1_command_proto_depIdxs,
		MessageInfos:      file_arrebato_message_command_v1_command_proto_msgTypes,
	}.Build()
	File_arrebato_message_command_v1_command_proto = out.File
	file_arrebato_message_command_v1_command_proto_rawDesc = nil
	file_arrebato_message_command_v1_command_proto_goTypes = nil
	file_arrebato_message_command_v1_command_proto_depIdxs = nil
}
