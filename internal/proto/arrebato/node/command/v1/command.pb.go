// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: arrebato/node/command/v1/command.proto

// Package arrebato.node.command.v1 provides commands written to the raft log regarding node state.
// These messages are used internally by nodes in the cluster to portray changes from the leader node to all followers.

package nodecmd

import (
	v1 "github.com/davidsbond/arrebato/internal/proto/arrebato/node/v1"
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

// The AddNode message is a command indicating that a node has been added to the cluster.
type AddNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *v1.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *AddNode) Reset() {
	*x = AddNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_node_command_v1_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddNode) ProtoMessage() {}

func (x *AddNode) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_node_command_v1_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddNode.ProtoReflect.Descriptor instead.
func (*AddNode) Descriptor() ([]byte, []int) {
	return file_arrebato_node_command_v1_command_proto_rawDescGZIP(), []int{0}
}

func (x *AddNode) GetNode() *v1.Node {
	if x != nil {
		return x.Node
	}
	return nil
}

// The RemoveNode message is a command indicating that a node has been removed from the cluster.
type RemoveNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *RemoveNode) Reset() {
	*x = RemoveNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_node_command_v1_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveNode) ProtoMessage() {}

func (x *RemoveNode) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_node_command_v1_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveNode.ProtoReflect.Descriptor instead.
func (*RemoveNode) Descriptor() ([]byte, []int) {
	return file_arrebato_node_command_v1_command_proto_rawDescGZIP(), []int{1}
}

func (x *RemoveNode) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The AssignTopic message is a command indicating that a topic has been assigned to a given node.
type AssignTopic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeName  string `protobuf:"bytes,1,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	TopicName string `protobuf:"bytes,2,opt,name=topic_name,json=topicName,proto3" json:"topic_name,omitempty"`
}

func (x *AssignTopic) Reset() {
	*x = AssignTopic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_node_command_v1_command_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignTopic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignTopic) ProtoMessage() {}

func (x *AssignTopic) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_node_command_v1_command_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignTopic.ProtoReflect.Descriptor instead.
func (*AssignTopic) Descriptor() ([]byte, []int) {
	return file_arrebato_node_command_v1_command_proto_rawDescGZIP(), []int{2}
}

func (x *AssignTopic) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *AssignTopic) GetTopicName() string {
	if x != nil {
		return x.TopicName
	}
	return ""
}

// The UnassignTopic message is a command indicating that a topic should be unassigned from all nodes. For example,
// if it is deleted.
type UnassignTopic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *UnassignTopic) Reset() {
	*x = UnassignTopic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_node_command_v1_command_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnassignTopic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnassignTopic) ProtoMessage() {}

func (x *UnassignTopic) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_node_command_v1_command_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnassignTopic.ProtoReflect.Descriptor instead.
func (*UnassignTopic) Descriptor() ([]byte, []int) {
	return file_arrebato_node_command_v1_command_proto_rawDescGZIP(), []int{3}
}

func (x *UnassignTopic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_arrebato_node_command_v1_command_proto protoreflect.FileDescriptor

var file_arrebato_node_command_v1_command_proto_rawDesc = []byte{
	0x0a, 0x26, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61,
	0x74, 0x6f, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e,
	0x76, 0x31, 0x1a, 0x1b, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x64,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x35, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x6e, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62,
	0x61, 0x74, 0x6f, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x20, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x49, 0x0a, 0x0b, 0x41, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x23, 0x0a, 0x0d, 0x55, 0x6e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x50, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76, 0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e,
	0x64, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61,
	0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f,
	0x76, 0x31, 0x3b, 0x6e, 0x6f, 0x64, 0x65, 0x63, 0x6d, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_arrebato_node_command_v1_command_proto_rawDescOnce sync.Once
	file_arrebato_node_command_v1_command_proto_rawDescData = file_arrebato_node_command_v1_command_proto_rawDesc
)

func file_arrebato_node_command_v1_command_proto_rawDescGZIP() []byte {
	file_arrebato_node_command_v1_command_proto_rawDescOnce.Do(func() {
		file_arrebato_node_command_v1_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_node_command_v1_command_proto_rawDescData)
	})
	return file_arrebato_node_command_v1_command_proto_rawDescData
}

var file_arrebato_node_command_v1_command_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_arrebato_node_command_v1_command_proto_goTypes = []interface{}{
	(*AddNode)(nil),       // 0: arrebato.node.command.v1.AddNode
	(*RemoveNode)(nil),    // 1: arrebato.node.command.v1.RemoveNode
	(*AssignTopic)(nil),   // 2: arrebato.node.command.v1.AssignTopic
	(*UnassignTopic)(nil), // 3: arrebato.node.command.v1.UnassignTopic
	(*v1.Node)(nil),       // 4: arrebato.node.v1.Node
}
var file_arrebato_node_command_v1_command_proto_depIdxs = []int32{
	4, // 0: arrebato.node.command.v1.AddNode.node:type_name -> arrebato.node.v1.Node
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_arrebato_node_command_v1_command_proto_init() }
func file_arrebato_node_command_v1_command_proto_init() {
	if File_arrebato_node_command_v1_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_node_command_v1_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddNode); i {
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
		file_arrebato_node_command_v1_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveNode); i {
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
		file_arrebato_node_command_v1_command_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignTopic); i {
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
		file_arrebato_node_command_v1_command_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnassignTopic); i {
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
			RawDescriptor: file_arrebato_node_command_v1_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_arrebato_node_command_v1_command_proto_goTypes,
		DependencyIndexes: file_arrebato_node_command_v1_command_proto_depIdxs,
		MessageInfos:      file_arrebato_node_command_v1_command_proto_msgTypes,
	}.Build()
	File_arrebato_node_command_v1_command_proto = out.File
	file_arrebato_node_command_v1_command_proto_rawDesc = nil
	file_arrebato_node_command_v1_command_proto_goTypes = nil
	file_arrebato_node_command_v1_command_proto_depIdxs = nil
}
