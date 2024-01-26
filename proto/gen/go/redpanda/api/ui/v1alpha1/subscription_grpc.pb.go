// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: redpanda/api/ui/v1alpha1/subscription.proto

package uiv1alpha1

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

const (
	SubscriptionService_GetCurrentSubscription_FullMethodName = "/redpanda.api.ui.v1alpha1.SubscriptionService/GetCurrentSubscription"
)

// SubscriptionServiceClient is the client API for SubscriptionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubscriptionServiceClient interface {
	GetCurrentSubscription(ctx context.Context, in *GetCurrentSubscriptionRequest, opts ...grpc.CallOption) (*GetCurrentSubscriptionResponse, error)
}

type subscriptionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSubscriptionServiceClient(cc grpc.ClientConnInterface) SubscriptionServiceClient {
	return &subscriptionServiceClient{cc}
}

func (c *subscriptionServiceClient) GetCurrentSubscription(ctx context.Context, in *GetCurrentSubscriptionRequest, opts ...grpc.CallOption) (*GetCurrentSubscriptionResponse, error) {
	out := new(GetCurrentSubscriptionResponse)
	err := c.cc.Invoke(ctx, SubscriptionService_GetCurrentSubscription_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubscriptionServiceServer is the server API for SubscriptionService service.
// All implementations should embed UnimplementedSubscriptionServiceServer
// for forward compatibility
type SubscriptionServiceServer interface {
	GetCurrentSubscription(context.Context, *GetCurrentSubscriptionRequest) (*GetCurrentSubscriptionResponse, error)
}

// UnimplementedSubscriptionServiceServer should be embedded to have forward compatible implementations.
type UnimplementedSubscriptionServiceServer struct {
}

func (UnimplementedSubscriptionServiceServer) GetCurrentSubscription(context.Context, *GetCurrentSubscriptionRequest) (*GetCurrentSubscriptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentSubscription not implemented")
}

// UnsafeSubscriptionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubscriptionServiceServer will
// result in compilation errors.
type UnsafeSubscriptionServiceServer interface {
	mustEmbedUnimplementedSubscriptionServiceServer()
}

func RegisterSubscriptionServiceServer(s grpc.ServiceRegistrar, srv SubscriptionServiceServer) {
	s.RegisterService(&SubscriptionService_ServiceDesc, srv)
}

func _SubscriptionService_GetCurrentSubscription_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentSubscriptionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriptionServiceServer).GetCurrentSubscription(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SubscriptionService_GetCurrentSubscription_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriptionServiceServer).GetCurrentSubscription(ctx, req.(*GetCurrentSubscriptionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SubscriptionService_ServiceDesc is the grpc.ServiceDesc for SubscriptionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SubscriptionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "redpanda.api.ui.v1alpha1.SubscriptionService",
	HandlerType: (*SubscriptionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentSubscription",
			Handler:    _SubscriptionService_GetCurrentSubscription_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "redpanda/api/ui/v1alpha1/subscription.proto",
}