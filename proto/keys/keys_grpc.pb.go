// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: keys/keys.proto

package keys

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KeysGreeterClient is the client API for KeysGreeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeysGreeterClient interface {
	UploadKeys(ctx context.Context, in *FileUploadRequest, opts ...grpc.CallOption) (*FileUploadResponse, error)
	DownloadKeys(ctx context.Context, in *Empty, opts ...grpc.CallOption) (KeysGreeter_DownloadKeysClient, error)
}

type keysGreeterClient struct {
	cc grpc.ClientConnInterface
}

func NewKeysGreeterClient(cc grpc.ClientConnInterface) KeysGreeterClient {
	return &keysGreeterClient{cc}
}

func (c *keysGreeterClient) UploadKeys(ctx context.Context, in *FileUploadRequest, opts ...grpc.CallOption) (*FileUploadResponse, error) {
	out := new(FileUploadResponse)
	err := c.cc.Invoke(ctx, "/keys.KeysGreeter/UploadKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keysGreeterClient) DownloadKeys(ctx context.Context, in *Empty, opts ...grpc.CallOption) (KeysGreeter_DownloadKeysClient, error) {
	stream, err := c.cc.NewStream(ctx, &KeysGreeter_ServiceDesc.Streams[0], "/keys.KeysGreeter/DownloadKeys", opts...)
	if err != nil {
		return nil, err
	}
	x := &keysGreeterDownloadKeysClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KeysGreeter_DownloadKeysClient interface {
	Recv() (*FileDownloadResponse, error)
	grpc.ClientStream
}

type keysGreeterDownloadKeysClient struct {
	grpc.ClientStream
}

func (x *keysGreeterDownloadKeysClient) Recv() (*FileDownloadResponse, error) {
	m := new(FileDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// KeysGreeterServer is the server API for KeysGreeter service.
// All implementations must embed UnimplementedKeysGreeterServer
// for forward compatibility
type KeysGreeterServer interface {
	UploadKeys(context.Context, *FileUploadRequest) (*FileUploadResponse, error)
	DownloadKeys(*Empty, KeysGreeter_DownloadKeysServer) error
	mustEmbedUnimplementedKeysGreeterServer()
}

// UnimplementedKeysGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedKeysGreeterServer struct {
}

func (UnimplementedKeysGreeterServer) UploadKeys(context.Context, *FileUploadRequest) (*FileUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadKeys not implemented")
}
func (UnimplementedKeysGreeterServer) DownloadKeys(*Empty, KeysGreeter_DownloadKeysServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadKeys not implemented")
}
func (UnimplementedKeysGreeterServer) mustEmbedUnimplementedKeysGreeterServer() {}

// UnsafeKeysGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeysGreeterServer will
// result in compilation errors.
type UnsafeKeysGreeterServer interface {
	mustEmbedUnimplementedKeysGreeterServer()
}

func RegisterKeysGreeterServer(s grpc.ServiceRegistrar, srv KeysGreeterServer) {
	s.RegisterService(&KeysGreeter_ServiceDesc, srv)
}

func _KeysGreeter_UploadKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeysGreeterServer).UploadKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/keys.KeysGreeter/UploadKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeysGreeterServer).UploadKeys(ctx, req.(*FileUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeysGreeter_DownloadKeys_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KeysGreeterServer).DownloadKeys(m, &keysGreeterDownloadKeysServer{stream})
}

type KeysGreeter_DownloadKeysServer interface {
	Send(*FileDownloadResponse) error
	grpc.ServerStream
}

type keysGreeterDownloadKeysServer struct {
	grpc.ServerStream
}

func (x *keysGreeterDownloadKeysServer) Send(m *FileDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

// KeysGreeter_ServiceDesc is the grpc.ServiceDesc for KeysGreeter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeysGreeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keys.KeysGreeter",
	HandlerType: (*KeysGreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadKeys",
			Handler:    _KeysGreeter_UploadKeys_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadKeys",
			Handler:       _KeysGreeter_DownloadKeys_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "keys/keys.proto",
}