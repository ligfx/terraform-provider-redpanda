// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: redpanda/api/billing/v1alpha1/invoice.proto

package billingv1alpha1

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
	InvoiceService_ListInvoices_FullMethodName = "/redpanda.api.billing.v1alpha1.InvoiceService/ListInvoices"
)

// InvoiceServiceClient is the client API for InvoiceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InvoiceServiceClient interface {
	ListInvoices(ctx context.Context, in *ListInvoicesRequest, opts ...grpc.CallOption) (*ListInvoicesResponse, error)
}

type invoiceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInvoiceServiceClient(cc grpc.ClientConnInterface) InvoiceServiceClient {
	return &invoiceServiceClient{cc}
}

func (c *invoiceServiceClient) ListInvoices(ctx context.Context, in *ListInvoicesRequest, opts ...grpc.CallOption) (*ListInvoicesResponse, error) {
	out := new(ListInvoicesResponse)
	err := c.cc.Invoke(ctx, InvoiceService_ListInvoices_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InvoiceServiceServer is the server API for InvoiceService service.
// All implementations should embed UnimplementedInvoiceServiceServer
// for forward compatibility
type InvoiceServiceServer interface {
	ListInvoices(context.Context, *ListInvoicesRequest) (*ListInvoicesResponse, error)
}

// UnimplementedInvoiceServiceServer should be embedded to have forward compatible implementations.
type UnimplementedInvoiceServiceServer struct {
}

func (UnimplementedInvoiceServiceServer) ListInvoices(context.Context, *ListInvoicesRequest) (*ListInvoicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInvoices not implemented")
}

// UnsafeInvoiceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InvoiceServiceServer will
// result in compilation errors.
type UnsafeInvoiceServiceServer interface {
	mustEmbedUnimplementedInvoiceServiceServer()
}

func RegisterInvoiceServiceServer(s grpc.ServiceRegistrar, srv InvoiceServiceServer) {
	s.RegisterService(&InvoiceService_ServiceDesc, srv)
}

func _InvoiceService_ListInvoices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInvoicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InvoiceServiceServer).ListInvoices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: InvoiceService_ListInvoices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InvoiceServiceServer).ListInvoices(ctx, req.(*ListInvoicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InvoiceService_ServiceDesc is the grpc.ServiceDesc for InvoiceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InvoiceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "redpanda.api.billing.v1alpha1.InvoiceService",
	HandlerType: (*InvoiceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListInvoices",
			Handler:    _InvoiceService_ListInvoices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "redpanda/api/billing/v1alpha1/invoice.proto",
}