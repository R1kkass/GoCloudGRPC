// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: access.proto

package access

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

type RequestAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FolderId int32 `protobuf:"varint,1,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	FileId   int32 `protobuf:"varint,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	UserId   int32 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RequestAccess) Reset() {
	*x = RequestAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAccess) ProtoMessage() {}

func (x *RequestAccess) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAccess.ProtoReflect.Descriptor instead.
func (*RequestAccess) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{0}
}

func (x *RequestAccess) GetFolderId() int32 {
	if x != nil {
		return x.FolderId
	}
	return 0
}

func (x *RequestAccess) GetFileId() int32 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *RequestAccess) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ResponseAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResponseAccess) Reset() {
	*x = ResponseAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseAccess) ProtoMessage() {}

func (x *ResponseAccess) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseAccess.ProtoReflect.Descriptor instead.
func (*ResponseAccess) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseAccess) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetAccessesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accesses []*RequestAccessData `protobuf:"bytes,1,rep,name=accesses,proto3" json:"accesses,omitempty"`
}

func (x *GetAccessesResponse) Reset() {
	*x = GetAccessesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccessesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessesResponse) ProtoMessage() {}

func (x *GetAccessesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessesResponse.ProtoReflect.Descriptor instead.
func (*GetAccessesResponse) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{2}
}

func (x *GetAccessesResponse) GetAccesses() []*RequestAccessData {
	if x != nil {
		return x.Accesses
	}
	return nil
}

type RequestAccessData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CurrentUserId int32 `protobuf:"varint,3,opt,name=current_user_id,json=currentUserId,proto3" json:"current_user_id,omitempty"`
	FileId        int32 `protobuf:"varint,4,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	FolderId      int32 `protobuf:"varint,5,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	StatusId      int32 `protobuf:"varint,6,opt,name=status_id,json=statusId,proto3" json:"status_id,omitempty"`
	User          *User `protobuf:"bytes,7,opt,name=user,proto3" json:"user,omitempty"`
	CurentUser    *User `protobuf:"bytes,8,opt,name=curent_user,json=curentUser,proto3" json:"curent_user,omitempty"`
}

func (x *RequestAccessData) Reset() {
	*x = RequestAccessData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestAccessData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestAccessData) ProtoMessage() {}

func (x *RequestAccessData) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestAccessData.ProtoReflect.Descriptor instead.
func (*RequestAccessData) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{3}
}

func (x *RequestAccessData) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RequestAccessData) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *RequestAccessData) GetCurrentUserId() int32 {
	if x != nil {
		return x.CurrentUserId
	}
	return 0
}

func (x *RequestAccessData) GetFileId() int32 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *RequestAccessData) GetFolderId() int32 {
	if x != nil {
		return x.FolderId
	}
	return 0
}

func (x *RequestAccessData) GetStatusId() int32 {
	if x != nil {
		return x.StatusId
	}
	return 0
}

func (x *RequestAccessData) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *RequestAccessData) GetCurentUser() *User {
	if x != nil {
		return x.CurentUser
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{4}
}

func (x *User) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{5}
}

type ChangeAccessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Status int32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ChangeAccessRequest) Reset() {
	*x = ChangeAccessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccessRequest) ProtoMessage() {}

func (x *ChangeAccessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccessRequest.ProtoReflect.Descriptor instead.
func (*ChangeAccessRequest) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{6}
}

func (x *ChangeAccessRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ChangeAccessRequest) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type ChangeAccessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ChangeAccessResponse) Reset() {
	*x = ChangeAccessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccessResponse) ProtoMessage() {}

func (x *ChangeAccessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_access_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccessResponse.ProtoReflect.Descriptor instead.
func (*ChangeAccessResponse) Descriptor() ([]byte, []int) {
	return file_access_proto_rawDescGZIP(), []int{7}
}

func (x *ChangeAccessResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_access_proto protoreflect.FileDescriptor

var file_access_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x5e, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2a, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x4c, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73,
	0x22, 0x88, 0x02, 0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x26, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x0b,
	0x63, 0x75, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x0a, 0x63, 0x75, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x22, 0x40, 0x0a, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x07, 0x0a,
	0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3d, 0x0a, 0x13, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x30, 0x0a, 0x14, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xda, 0x01, 0x0a, 0x0d, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x15, 0x2e, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x1a, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0b, 0x47, 0x65,
	0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x0d, 0x2e, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_access_proto_rawDescOnce sync.Once
	file_access_proto_rawDescData = file_access_proto_rawDesc
)

func file_access_proto_rawDescGZIP() []byte {
	file_access_proto_rawDescOnce.Do(func() {
		file_access_proto_rawDescData = protoimpl.X.CompressGZIP(file_access_proto_rawDescData)
	})
	return file_access_proto_rawDescData
}

var file_access_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_access_proto_goTypes = []interface{}{
	(*RequestAccess)(nil),        // 0: access.RequestAccess
	(*ResponseAccess)(nil),       // 1: access.ResponseAccess
	(*GetAccessesResponse)(nil),  // 2: access.GetAccessesResponse
	(*RequestAccessData)(nil),    // 3: access.RequestAccessData
	(*User)(nil),                 // 4: access.User
	(*Empty)(nil),                // 5: access.Empty
	(*ChangeAccessRequest)(nil),  // 6: access.ChangeAccessRequest
	(*ChangeAccessResponse)(nil), // 7: access.ChangeAccessResponse
}
var file_access_proto_depIdxs = []int32{
	3, // 0: access.GetAccessesResponse.accesses:type_name -> access.RequestAccessData
	4, // 1: access.RequestAccessData.user:type_name -> access.User
	4, // 2: access.RequestAccessData.curent_user:type_name -> access.User
	0, // 3: access.AccessGreeter.CreateAccess:input_type -> access.RequestAccess
	5, // 4: access.AccessGreeter.GetAccesses:input_type -> access.Empty
	6, // 5: access.AccessGreeter.ChangeAccess:input_type -> access.ChangeAccessRequest
	1, // 6: access.AccessGreeter.CreateAccess:output_type -> access.ResponseAccess
	2, // 7: access.AccessGreeter.GetAccesses:output_type -> access.GetAccessesResponse
	7, // 8: access.AccessGreeter.ChangeAccess:output_type -> access.ChangeAccessResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_access_proto_init() }
func file_access_proto_init() {
	if File_access_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_access_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestAccess); i {
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
		file_access_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseAccess); i {
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
		file_access_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccessesResponse); i {
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
		file_access_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestAccessData); i {
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
		file_access_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_access_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_access_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccessRequest); i {
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
		file_access_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccessResponse); i {
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
			RawDescriptor: file_access_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_access_proto_goTypes,
		DependencyIndexes: file_access_proto_depIdxs,
		MessageInfos:      file_access_proto_msgTypes,
	}.Build()
	File_access_proto = out.File
	file_access_proto_rawDesc = nil
	file_access_proto_goTypes = nil
	file_access_proto_depIdxs = nil
}
