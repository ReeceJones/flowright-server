// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: routing.proto

package flowright

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

// RoutingControllerClient is the client API for RoutingController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoutingControllerClient interface {
	RegisterProxy(ctx context.Context, in *ProxyRegisterRequest, opts ...grpc.CallOption) (*ProxyRegisterResponse, error)
	GetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*ControllerInfo, error)
	CreateOrUpdateRoute(ctx context.Context, in *RoutingMap, opts ...grpc.CallOption) (*RoutingRule, error)
	GetRoute(ctx context.Context, in *RoutingParams, opts ...grpc.CallOption) (*RoutingRule, error)
	// TODO
	CreateEnvironment(ctx context.Context, in *EnvironmentCreateRequest, opts ...grpc.CallOption) (*EnvironmentCreateResponse, error)
	DeleteEnvironment(ctx context.Context, in *EnvironmentDeleteRequest, opts ...grpc.CallOption) (*EnvironmentDeleteResponse, error)
	ListEnvironments(ctx context.Context, in *EnvironmentListRequest, opts ...grpc.CallOption) (*EnvironmentListResponse, error)
}

type routingControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewRoutingControllerClient(cc grpc.ClientConnInterface) RoutingControllerClient {
	return &routingControllerClient{cc}
}

func (c *routingControllerClient) RegisterProxy(ctx context.Context, in *ProxyRegisterRequest, opts ...grpc.CallOption) (*ProxyRegisterResponse, error) {
	out := new(ProxyRegisterResponse)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/RegisterProxy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) GetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*ControllerInfo, error) {
	out := new(ControllerInfo)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) CreateOrUpdateRoute(ctx context.Context, in *RoutingMap, opts ...grpc.CallOption) (*RoutingRule, error) {
	out := new(RoutingRule)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/CreateOrUpdateRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) GetRoute(ctx context.Context, in *RoutingParams, opts ...grpc.CallOption) (*RoutingRule, error) {
	out := new(RoutingRule)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/GetRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) CreateEnvironment(ctx context.Context, in *EnvironmentCreateRequest, opts ...grpc.CallOption) (*EnvironmentCreateResponse, error) {
	out := new(EnvironmentCreateResponse)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/CreateEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) DeleteEnvironment(ctx context.Context, in *EnvironmentDeleteRequest, opts ...grpc.CallOption) (*EnvironmentDeleteResponse, error) {
	out := new(EnvironmentDeleteResponse)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/DeleteEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingControllerClient) ListEnvironments(ctx context.Context, in *EnvironmentListRequest, opts ...grpc.CallOption) (*EnvironmentListResponse, error) {
	out := new(EnvironmentListResponse)
	err := c.cc.Invoke(ctx, "/routing.RoutingController/ListEnvironments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoutingControllerServer is the server API for RoutingController service.
// All implementations must embed UnimplementedRoutingControllerServer
// for forward compatibility
type RoutingControllerServer interface {
	RegisterProxy(context.Context, *ProxyRegisterRequest) (*ProxyRegisterResponse, error)
	GetInfo(context.Context, *InfoRequest) (*ControllerInfo, error)
	CreateOrUpdateRoute(context.Context, *RoutingMap) (*RoutingRule, error)
	GetRoute(context.Context, *RoutingParams) (*RoutingRule, error)
	// TODO
	CreateEnvironment(context.Context, *EnvironmentCreateRequest) (*EnvironmentCreateResponse, error)
	DeleteEnvironment(context.Context, *EnvironmentDeleteRequest) (*EnvironmentDeleteResponse, error)
	ListEnvironments(context.Context, *EnvironmentListRequest) (*EnvironmentListResponse, error)
	mustEmbedUnimplementedRoutingControllerServer()
}

// UnimplementedRoutingControllerServer must be embedded to have forward compatible implementations.
type UnimplementedRoutingControllerServer struct {
}

func (UnimplementedRoutingControllerServer) RegisterProxy(context.Context, *ProxyRegisterRequest) (*ProxyRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterProxy not implemented")
}
func (UnimplementedRoutingControllerServer) GetInfo(context.Context, *InfoRequest) (*ControllerInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedRoutingControllerServer) CreateOrUpdateRoute(context.Context, *RoutingMap) (*RoutingRule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateRoute not implemented")
}
func (UnimplementedRoutingControllerServer) GetRoute(context.Context, *RoutingParams) (*RoutingRule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoute not implemented")
}
func (UnimplementedRoutingControllerServer) CreateEnvironment(context.Context, *EnvironmentCreateRequest) (*EnvironmentCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEnvironment not implemented")
}
func (UnimplementedRoutingControllerServer) DeleteEnvironment(context.Context, *EnvironmentDeleteRequest) (*EnvironmentDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEnvironment not implemented")
}
func (UnimplementedRoutingControllerServer) ListEnvironments(context.Context, *EnvironmentListRequest) (*EnvironmentListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEnvironments not implemented")
}
func (UnimplementedRoutingControllerServer) mustEmbedUnimplementedRoutingControllerServer() {}

// UnsafeRoutingControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoutingControllerServer will
// result in compilation errors.
type UnsafeRoutingControllerServer interface {
	mustEmbedUnimplementedRoutingControllerServer()
}

func RegisterRoutingControllerServer(s grpc.ServiceRegistrar, srv RoutingControllerServer) {
	s.RegisterService(&RoutingController_ServiceDesc, srv)
}

func _RoutingController_RegisterProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProxyRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).RegisterProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/RegisterProxy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).RegisterProxy(ctx, req.(*ProxyRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).GetInfo(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_CreateOrUpdateRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoutingMap)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).CreateOrUpdateRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/CreateOrUpdateRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).CreateOrUpdateRoute(ctx, req.(*RoutingMap))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_GetRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoutingParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).GetRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/GetRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).GetRoute(ctx, req.(*RoutingParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_CreateEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnvironmentCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).CreateEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/CreateEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).CreateEnvironment(ctx, req.(*EnvironmentCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_DeleteEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnvironmentDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).DeleteEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/DeleteEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).DeleteEnvironment(ctx, req.(*EnvironmentDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingController_ListEnvironments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnvironmentListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingControllerServer).ListEnvironments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/routing.RoutingController/ListEnvironments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingControllerServer).ListEnvironments(ctx, req.(*EnvironmentListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoutingController_ServiceDesc is the grpc.ServiceDesc for RoutingController service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoutingController_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "routing.RoutingController",
	HandlerType: (*RoutingControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterProxy",
			Handler:    _RoutingController_RegisterProxy_Handler,
		},
		{
			MethodName: "GetInfo",
			Handler:    _RoutingController_GetInfo_Handler,
		},
		{
			MethodName: "CreateOrUpdateRoute",
			Handler:    _RoutingController_CreateOrUpdateRoute_Handler,
		},
		{
			MethodName: "GetRoute",
			Handler:    _RoutingController_GetRoute_Handler,
		},
		{
			MethodName: "CreateEnvironment",
			Handler:    _RoutingController_CreateEnvironment_Handler,
		},
		{
			MethodName: "DeleteEnvironment",
			Handler:    _RoutingController_DeleteEnvironment_Handler,
		},
		{
			MethodName: "ListEnvironments",
			Handler:    _RoutingController_ListEnvironments_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routing.proto",
}