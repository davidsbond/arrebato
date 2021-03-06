// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: arrebato/signing/service/v1/service.proto

// Package arrebato.signing.service.v1 provides the schema for the SigningService, which is used to manage signing keys within the
// cluster.

package signingsvc

import (
	v1 "github.com/davidsbond/arrebato/internal/proto/arrebato/signing/v1"
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

// The CreateKeyPairRequest message is the request DTO when calling SigningService.CreateKeyPair.
type CreateKeyPairRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateKeyPairRequest) Reset() {
	*x = CreateKeyPairRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateKeyPairRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateKeyPairRequest) ProtoMessage() {}

func (x *CreateKeyPairRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateKeyPairRequest.ProtoReflect.Descriptor instead.
func (*CreateKeyPairRequest) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{0}
}

// The CreateKeyPairResponse message is the response DTO when calling SigningService.CreateKeyPair.
type CreateKeyPairResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The generated public/private key pair.
	KeyPair *v1.KeyPair `protobuf:"bytes,1,opt,name=key_pair,json=keyPair,proto3" json:"key_pair,omitempty"`
}

func (x *CreateKeyPairResponse) Reset() {
	*x = CreateKeyPairResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateKeyPairResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateKeyPairResponse) ProtoMessage() {}

func (x *CreateKeyPairResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateKeyPairResponse.ProtoReflect.Descriptor instead.
func (*CreateKeyPairResponse) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateKeyPairResponse) GetKeyPair() *v1.KeyPair {
	if x != nil {
		return x.KeyPair
	}
	return nil
}

// The GetPublicKeyRequest message is the request DTO when calling SigningService.GetPublicKey.
type GetPublicKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The client to obtain the public key of.
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
}

func (x *GetPublicKeyRequest) Reset() {
	*x = GetPublicKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPublicKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPublicKeyRequest) ProtoMessage() {}

func (x *GetPublicKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPublicKeyRequest.ProtoReflect.Descriptor instead.
func (*GetPublicKeyRequest) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetPublicKeyRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

// The GetPublicKeyResponse message is the request DTO when calling SigningService.GetPublicKey.
type GetPublicKeyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The client's public key.
	PublicKey *v1.PublicKey `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *GetPublicKeyResponse) Reset() {
	*x = GetPublicKeyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPublicKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPublicKeyResponse) ProtoMessage() {}

func (x *GetPublicKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPublicKeyResponse.ProtoReflect.Descriptor instead.
func (*GetPublicKeyResponse) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetPublicKeyResponse) GetPublicKey() *v1.PublicKey {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

// The ListPublicKeysRequest message is the request DTO when calling SigningService.ListPublicKeys.
type ListPublicKeysRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListPublicKeysRequest) Reset() {
	*x = ListPublicKeysRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPublicKeysRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPublicKeysRequest) ProtoMessage() {}

func (x *ListPublicKeysRequest) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPublicKeysRequest.ProtoReflect.Descriptor instead.
func (*ListPublicKeysRequest) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{4}
}

// The ListPublicKeysRequest message is the response DTO when calling SigningService.ListPublicKeys.
type ListPublicKeysResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The public keys stored in the server.
	PublicKeys []*v1.PublicKey `protobuf:"bytes,1,rep,name=public_keys,json=publicKeys,proto3" json:"public_keys,omitempty"`
}

func (x *ListPublicKeysResponse) Reset() {
	*x = ListPublicKeysResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPublicKeysResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPublicKeysResponse) ProtoMessage() {}

func (x *ListPublicKeysResponse) ProtoReflect() protoreflect.Message {
	mi := &file_arrebato_signing_service_v1_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPublicKeysResponse.ProtoReflect.Descriptor instead.
func (*ListPublicKeysResponse) Descriptor() ([]byte, []int) {
	return file_arrebato_signing_service_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListPublicKeysResponse) GetPublicKeys() []*v1.PublicKey {
	if x != nil {
		return x.PublicKeys
	}
	return nil
}

var File_arrebato_signing_service_v1_service_proto protoreflect.FileDescriptor

var file_arrebato_signing_service_v1_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x61, 0x72, 0x72,
	0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61,
	0x74, 0x6f, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69,
	0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x16, 0x0a, 0x14, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x50, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79,
	0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x08,
	0x6b, 0x65, 0x79, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x69, 0x72, 0x52, 0x07, 0x6b, 0x65,
	0x79, 0x50, 0x61, 0x69, 0x72, 0x22, 0x32, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x55, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f,
	0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x22, 0x17, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x59, 0x0a, 0x16, 0x4c, 0x69, 0x73,
	0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0b, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65,
	0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62,
	0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x73, 0x32, 0xf8, 0x02, 0x0a, 0x0e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x76, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x69, 0x72, 0x12, 0x31, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62,
	0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79,
	0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x61, 0x72,
	0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4b, 0x65, 0x79, 0x50, 0x61, 0x69, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x73, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x30, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69,
	0x6e, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x31, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x79, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x32, 0x2e, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74,
	0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x61, 0x72, 0x72,
	0x65, 0x62, 0x61, 0x74, 0x6f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x56, 0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61,
	0x76, 0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e, 0x64, 0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74,
	0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x61, 0x72, 0x72, 0x65, 0x62, 0x61, 0x74, 0x6f, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e,
	0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x73, 0x76, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_arrebato_signing_service_v1_service_proto_rawDescOnce sync.Once
	file_arrebato_signing_service_v1_service_proto_rawDescData = file_arrebato_signing_service_v1_service_proto_rawDesc
)

func file_arrebato_signing_service_v1_service_proto_rawDescGZIP() []byte {
	file_arrebato_signing_service_v1_service_proto_rawDescOnce.Do(func() {
		file_arrebato_signing_service_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_arrebato_signing_service_v1_service_proto_rawDescData)
	})
	return file_arrebato_signing_service_v1_service_proto_rawDescData
}

var file_arrebato_signing_service_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_arrebato_signing_service_v1_service_proto_goTypes = []interface{}{
	(*CreateKeyPairRequest)(nil),   // 0: arrebato.signing.service.v1.CreateKeyPairRequest
	(*CreateKeyPairResponse)(nil),  // 1: arrebato.signing.service.v1.CreateKeyPairResponse
	(*GetPublicKeyRequest)(nil),    // 2: arrebato.signing.service.v1.GetPublicKeyRequest
	(*GetPublicKeyResponse)(nil),   // 3: arrebato.signing.service.v1.GetPublicKeyResponse
	(*ListPublicKeysRequest)(nil),  // 4: arrebato.signing.service.v1.ListPublicKeysRequest
	(*ListPublicKeysResponse)(nil), // 5: arrebato.signing.service.v1.ListPublicKeysResponse
	(*v1.KeyPair)(nil),             // 6: arrebato.signing.v1.KeyPair
	(*v1.PublicKey)(nil),           // 7: arrebato.signing.v1.PublicKey
}
var file_arrebato_signing_service_v1_service_proto_depIdxs = []int32{
	6, // 0: arrebato.signing.service.v1.CreateKeyPairResponse.key_pair:type_name -> arrebato.signing.v1.KeyPair
	7, // 1: arrebato.signing.service.v1.GetPublicKeyResponse.public_key:type_name -> arrebato.signing.v1.PublicKey
	7, // 2: arrebato.signing.service.v1.ListPublicKeysResponse.public_keys:type_name -> arrebato.signing.v1.PublicKey
	0, // 3: arrebato.signing.service.v1.SigningService.CreateKeyPair:input_type -> arrebato.signing.service.v1.CreateKeyPairRequest
	2, // 4: arrebato.signing.service.v1.SigningService.GetPublicKey:input_type -> arrebato.signing.service.v1.GetPublicKeyRequest
	4, // 5: arrebato.signing.service.v1.SigningService.ListPublicKeys:input_type -> arrebato.signing.service.v1.ListPublicKeysRequest
	1, // 6: arrebato.signing.service.v1.SigningService.CreateKeyPair:output_type -> arrebato.signing.service.v1.CreateKeyPairResponse
	3, // 7: arrebato.signing.service.v1.SigningService.GetPublicKey:output_type -> arrebato.signing.service.v1.GetPublicKeyResponse
	5, // 8: arrebato.signing.service.v1.SigningService.ListPublicKeys:output_type -> arrebato.signing.service.v1.ListPublicKeysResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_arrebato_signing_service_v1_service_proto_init() }
func file_arrebato_signing_service_v1_service_proto_init() {
	if File_arrebato_signing_service_v1_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_arrebato_signing_service_v1_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateKeyPairRequest); i {
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
		file_arrebato_signing_service_v1_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateKeyPairResponse); i {
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
		file_arrebato_signing_service_v1_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPublicKeyRequest); i {
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
		file_arrebato_signing_service_v1_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPublicKeyResponse); i {
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
		file_arrebato_signing_service_v1_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPublicKeysRequest); i {
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
		file_arrebato_signing_service_v1_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPublicKeysResponse); i {
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
			RawDescriptor: file_arrebato_signing_service_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_arrebato_signing_service_v1_service_proto_goTypes,
		DependencyIndexes: file_arrebato_signing_service_v1_service_proto_depIdxs,
		MessageInfos:      file_arrebato_signing_service_v1_service_proto_msgTypes,
	}.Build()
	File_arrebato_signing_service_v1_service_proto = out.File
	file_arrebato_signing_service_v1_service_proto_rawDesc = nil
	file_arrebato_signing_service_v1_service_proto_goTypes = nil
	file_arrebato_signing_service_v1_service_proto_depIdxs = nil
}
