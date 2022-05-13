// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package go_grpc

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

// TestSvcClient is the client API for TestSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestSvcClient interface {
	Method1(ctx context.Context, in *Method1Req, opts ...grpc.CallOption) (*Method1Resp, error)
	Method2(ctx context.Context, in *Method2Req, opts ...grpc.CallOption) (*Method2Resp, error)
	Method3(ctx context.Context, in *Method2Req, opts ...grpc.CallOption) (*Method2Resp, error)
}

type testSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewTestSvcClient(cc grpc.ClientConnInterface) TestSvcClient {
	return &testSvcClient{cc}
}

func (c *testSvcClient) Method1(ctx context.Context, in *Method1Req, opts ...grpc.CallOption) (*Method1Resp, error) {
	out := new(Method1Resp)
	err := c.cc.Invoke(ctx, "/TestSvc.TestSvc/Method1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testSvcClient) Method2(ctx context.Context, in *Method2Req, opts ...grpc.CallOption) (*Method2Resp, error) {
	out := new(Method2Resp)
	err := c.cc.Invoke(ctx, "/TestSvc.TestSvc/Method2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testSvcClient) Method3(ctx context.Context, in *Method2Req, opts ...grpc.CallOption) (*Method2Resp, error) {
	out := new(Method2Resp)
	err := c.cc.Invoke(ctx, "/TestSvc.TestSvc/Method3", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestSvcServer is the server API for TestSvc service.
// All implementations should embed UnimplementedTestSvcServer
// for forward compatibility
type TestSvcServer interface {
	Method1(context.Context, *Method1Req) (*Method1Resp, error)
	Method2(context.Context, *Method2Req) (*Method2Resp, error)
	Method3(context.Context, *Method2Req) (*Method2Resp, error)
}

// UnimplementedTestSvcServer should be embedded to have forward compatible implementations.
type UnimplementedTestSvcServer struct {
}

func (UnimplementedTestSvcServer) Method1(context.Context, *Method1Req) (*Method1Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Method1 not implemented")
}
func (UnimplementedTestSvcServer) Method2(context.Context, *Method2Req) (*Method2Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Method2 not implemented")
}
func (UnimplementedTestSvcServer) Method3(context.Context, *Method2Req) (*Method2Resp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Method3 not implemented")
}

// UnsafeTestSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestSvcServer will
// result in compilation errors.
type UnsafeTestSvcServer interface {
	mustEmbedUnimplementedTestSvcServer()
}

func RegisterTestSvcServer(s grpc.ServiceRegistrar, srv TestSvcServer) {
	s.RegisterService(&TestSvc_ServiceDesc, srv)
}

func _TestSvc_Method1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Method1Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestSvcServer).Method1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TestSvc.TestSvc/Method1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestSvcServer).Method1(ctx, req.(*Method1Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestSvc_Method2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Method2Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestSvcServer).Method2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TestSvc.TestSvc/Method2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestSvcServer).Method2(ctx, req.(*Method2Req))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestSvc_Method3_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Method2Req)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestSvcServer).Method3(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TestSvc.TestSvc/Method3",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestSvcServer).Method3(ctx, req.(*Method2Req))
	}
	return interceptor(ctx, in, info, handler)
}

// TestSvc_ServiceDesc is the grpc.ServiceDesc for TestSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TestSvc.TestSvc",
	HandlerType: (*TestSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Method1",
			Handler:    _TestSvc_Method1_Handler,
		},
		{
			MethodName: "Method2",
			Handler:    _TestSvc_Method2_Handler,
		},
		{
			MethodName: "Method3",
			Handler:    _TestSvc_Method3_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/test.proto",
}
