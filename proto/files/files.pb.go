// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: files/files.proto

package files

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

type FileUploadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk    []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=fileName,proto3" json:"fileName,omitempty"`
	FolderId uint32 `protobuf:"varint,3,opt,name=folderId,proto3" json:"folderId,omitempty"`
}

func (x *FileUploadRequest) Reset() {
	*x = FileUploadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_files_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadRequest) ProtoMessage() {}

func (x *FileUploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_files_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadRequest.ProtoReflect.Descriptor instead.
func (*FileUploadRequest) Descriptor() ([]byte, []int) {
	return file_files_files_proto_rawDescGZIP(), []int{0}
}

func (x *FileUploadRequest) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *FileUploadRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileUploadRequest) GetFolderId() uint32 {
	if x != nil {
		return x.FolderId
	}
	return 0
}

type FileUploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *FileUploadResponse) Reset() {
	*x = FileUploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_files_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadResponse) ProtoMessage() {}

func (x *FileUploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_files_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadResponse.ProtoReflect.Descriptor instead.
func (*FileUploadResponse) Descriptor() ([]byte, []int) {
	return file_files_files_proto_rawDescGZIP(), []int{1}
}

func (x *FileUploadResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type FileDownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId   uint32 `protobuf:"varint,1,opt,name=fileId,proto3" json:"fileId,omitempty"`
	FolderId uint32 `protobuf:"varint,2,opt,name=folderId,proto3" json:"folderId,omitempty"`
}

func (x *FileDownloadRequest) Reset() {
	*x = FileDownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_files_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadRequest) ProtoMessage() {}

func (x *FileDownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_files_files_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadRequest.ProtoReflect.Descriptor instead.
func (*FileDownloadRequest) Descriptor() ([]byte, []int) {
	return file_files_files_proto_rawDescGZIP(), []int{2}
}

func (x *FileDownloadRequest) GetFileId() uint32 {
	if x != nil {
		return x.FileId
	}
	return 0
}

func (x *FileDownloadRequest) GetFolderId() uint32 {
	if x != nil {
		return x.FolderId
	}
	return 0
}

type FileDownloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk    []byte  `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	FileName string  `protobuf:"bytes,2,opt,name=fileName,proto3" json:"fileName,omitempty"`
	Progress float32 `protobuf:"fixed32,3,opt,name=progress,proto3" json:"progress,omitempty"`
}

func (x *FileDownloadResponse) Reset() {
	*x = FileDownloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_files_files_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDownloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDownloadResponse) ProtoMessage() {}

func (x *FileDownloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_files_files_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDownloadResponse.ProtoReflect.Descriptor instead.
func (*FileDownloadResponse) Descriptor() ([]byte, []int) {
	return file_files_files_proto_rawDescGZIP(), []int{3}
}

func (x *FileDownloadResponse) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *FileDownloadResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileDownloadResponse) GetProgress() float32 {
	if x != nil {
		return x.Progress
	}
	return 0
}

var File_files_files_proto protoreflect.FileDescriptor

var file_files_files_proto_rawDesc = []byte{
	0x0a, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x61, 0x0a, 0x11, 0x46, 0x69,
	0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2e, 0x0a,
	0x12, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x49, 0x0a,
	0x13, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x66, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x64, 0x0a, 0x14, 0x46, 0x69, 0x6c, 0x65,
	0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x32, 0x9e,
	0x01, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12,
	0x43, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x18, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x12, 0x49, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42,
	0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_files_files_proto_rawDescOnce sync.Once
	file_files_files_proto_rawDescData = file_files_files_proto_rawDesc
)

func file_files_files_proto_rawDescGZIP() []byte {
	file_files_files_proto_rawDescOnce.Do(func() {
		file_files_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_files_files_proto_rawDescData)
	})
	return file_files_files_proto_rawDescData
}

var file_files_files_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_files_files_proto_goTypes = []interface{}{
	(*FileUploadRequest)(nil),    // 0: files.FileUploadRequest
	(*FileUploadResponse)(nil),   // 1: files.FileUploadResponse
	(*FileDownloadRequest)(nil),  // 2: files.FileDownloadRequest
	(*FileDownloadResponse)(nil), // 3: files.FileDownloadResponse
}
var file_files_files_proto_depIdxs = []int32{
	0, // 0: files.FilesGreeter.UploadFile:input_type -> files.FileUploadRequest
	2, // 1: files.FilesGreeter.DownloadFile:input_type -> files.FileDownloadRequest
	1, // 2: files.FilesGreeter.UploadFile:output_type -> files.FileUploadResponse
	3, // 3: files.FilesGreeter.DownloadFile:output_type -> files.FileDownloadResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_files_files_proto_init() }
func file_files_files_proto_init() {
	if File_files_files_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_files_files_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadRequest); i {
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
		file_files_files_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadResponse); i {
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
		file_files_files_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadRequest); i {
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
		file_files_files_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDownloadResponse); i {
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
			RawDescriptor: file_files_files_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_files_files_proto_goTypes,
		DependencyIndexes: file_files_files_proto_depIdxs,
		MessageInfos:      file_files_files_proto_msgTypes,
	}.Build()
	File_files_files_proto = out.File
	file_files_files_proto_rawDesc = nil
	file_files_files_proto_goTypes = nil
	file_files_files_proto_depIdxs = nil
}