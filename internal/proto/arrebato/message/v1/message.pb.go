// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: arrebato/message/v1/message.proto

// Package arrebato.message.v1 provides domain models for messages.

package message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The Message message describes a single message stored within a topic.
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The topic the message belongs to.
	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	// The location of the message within the topic. This is managed by the cluster and should not be provided when
	// producing messages.
	Index uint64 `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	// The client-provided message key contents.
	Key *anypb.Any `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	// The client-provided message contents.
	Value *anypb.Any `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	// The time at which the message was produced. This is managed by the cluster and should not be provided when
	// producing messages.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Information about the client that produced the message.
	Sender *Sender `protobuf:"bytes,6,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_message_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_message_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_arrebato_message_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Message) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Message) GetKey() *anypb.Any {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Message) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Message) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Message) GetSender() *Sender {
	if x != nil {
		return x.Sender
	}
	return nil
}

// The Sender message describes a client that produced a Message.
type Sender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The client identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Whether or not the client provided the signed message key which has been verified by the server. This is used
	// to confirm the sender is indeed who they claim to be.
	Verified bool `protobuf:"varint,2,opt,name=verified,proto3" json:"verified,omitempty"`
	// The client-generated signature of the message key.
	KeySignature []byte `protobuf:"bytes,3,opt,name=key_signature,json=keySignature,proto3" json:"key_signature,omitempty"`
}

func (x *Sender) Reset() {
	*x = Sender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_message_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sender) ProtoMessage() {}

func (x *Sender) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_message_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sender.ProtoReflect.Descriptor instead.
func (*Sender) Descriptor() ([]byte, []int) {
	return file_arrebato_message_v1_message_proto_rawDescGZIP(), []int{1}
}

func (x *Sender) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Sender) GetVerified() bool {
	if x != nil {
		return x.Verified
	}
	return false
}

func (x *Sender) GetKeySignature() []byte {
	if x != nil {
		return x.KeySignature
	}
	return nil
}

var File_arrebato_message_v1_message_proto protoreflect.FileDescriptor

var file_arrebato_message_v1_message_proto_rawDesc = []byte{
	0x0a, 0x21, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf8, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x26, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x33, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x72, 0x72,
	0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x22,
	0x59, 0x0a, 0x06, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x6b, 0x65, 0x79, 0x5f, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x6b, 0x65,
	0x79, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76, 0x69, 0x64, 0x73, 0x62,
	0x6f, 0x6e, 0x64, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x72, 0x65,
	0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x3b,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arrebato_message_v1_message_proto_rawDescOnce sync.Once
	file_arrebato_message_v1_message_proto_rawDescData = file_arrebato_message_v1_message_proto_rawDesc
)

func file_arrebato_message_v1_message_proto_rawDescGZIP() []byte {
	file_arrebato_message_v1_message_proto_rawDescOnce.Do(func() {
		file_arrebato_message_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_message_v1_message_proto_rawDescData)
	})
	return file_arrebato_message_v1_message_proto_rawDescData
}

var file_arrebato_message_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_arrebato_message_v1_message_proto_goTypes = []interface{}{
	(*Message)(nil),               // 0: arrebato.message.v1.Message
	(*Sender)(nil),                // 1: arrebato.message.v1.Sender
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_arrebato_message_v1_message_proto_depIdxs = []int32{
	2, // 0: arrebato.message.v1.Message.key:type_name -> google.protobuf.Any
	2, // 1: arrebato.message.v1.Message.value:type_name -> google.protobuf.Any
	3, // 2: arrebato.message.v1.Message.timestamp:type_name -> google.protobuf.Timestamp
	1, // 3: arrebato.message.v1.Message.sender:type_name -> arrebato.message.v1.Sender
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_arrebato_message_v1_message_proto_init() }
func file_arrebato_message_v1_message_proto_init() {
	if File_arrebato_message_v1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_message_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_arrebato_message_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sender); i {
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
			RawDescriptor: file_arrebato_message_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arrebato_message_v1_message_proto_goTypes,
		DependencyIndexes: file_arrebato_message_v1_message_proto_depIdxs,
		MessageInfos:      file_arrebato_message_v1_message_proto_msgTypes,
	}.Build()
	File_arrebato_message_v1_message_proto = out.File
	file_arrebato_message_v1_message_proto_rawDesc = nil
	file_arrebato_message_v1_message_proto_goTypes = nil
	file_arrebato_message_v1_message_proto_depIdxs = nil
}
