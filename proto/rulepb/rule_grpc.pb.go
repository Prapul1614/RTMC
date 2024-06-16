// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: rule.proto

package rulepb

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

// RuleServiceClient is the client API for RuleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RuleServiceClient interface {
	CreateRule(ctx context.Context, in *CreateRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error)
	GetRules(ctx context.Context, in *GetRulesRequest, opts ...grpc.CallOption) (*RulesResponse, error)
	Classify(ctx context.Context, in *ClassifyRequest, opts ...grpc.CallOption) (*ClassifyResponse, error)
}

type ruleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRuleServiceClient(cc grpc.ClientConnInterface) RuleServiceClient {
	return &ruleServiceClient{cc}
}

func (c *ruleServiceClient) CreateRule(ctx context.Context, in *CreateRuleRequest, opts ...grpc.CallOption) (*RuleResponse, error) {
	out := new(RuleResponse)
	err := c.cc.Invoke(ctx, "/rulepb.RuleService/CreateRule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleServiceClient) GetRules(ctx context.Context, in *GetRulesRequest, opts ...grpc.CallOption) (*RulesResponse, error) {
	out := new(RulesResponse)
	err := c.cc.Invoke(ctx, "/rulepb.RuleService/GetRules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleServiceClient) Classify(ctx context.Context, in *ClassifyRequest, opts ...grpc.CallOption) (*ClassifyResponse, error) {
	out := new(ClassifyResponse)
	err := c.cc.Invoke(ctx, "/rulepb.RuleService/Classify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuleServiceServer is the server API for RuleService service.
// All implementations must embed UnimplementedRuleServiceServer
// for forward compatibility
type RuleServiceServer interface {
	CreateRule(context.Context, *CreateRuleRequest) (*RuleResponse, error)
	GetRules(context.Context, *GetRulesRequest) (*RulesResponse, error)
	Classify(context.Context, *ClassifyRequest) (*ClassifyResponse, error)
	mustEmbedUnimplementedRuleServiceServer()
}

// UnimplementedRuleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRuleServiceServer struct {
}

func (UnimplementedRuleServiceServer) CreateRule(context.Context, *CreateRuleRequest) (*RuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRule not implemented")
}
func (UnimplementedRuleServiceServer) GetRules(context.Context, *GetRulesRequest) (*RulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRules not implemented")
}
func (UnimplementedRuleServiceServer) Classify(context.Context, *ClassifyRequest) (*ClassifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Classify not implemented")
}
func (UnimplementedRuleServiceServer) mustEmbedUnimplementedRuleServiceServer() {}

// UnsafeRuleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RuleServiceServer will
// result in compilation errors.
type UnsafeRuleServiceServer interface {
	mustEmbedUnimplementedRuleServiceServer()
}

func RegisterRuleServiceServer(s grpc.ServiceRegistrar, srv RuleServiceServer) {
	s.RegisterService(&RuleService_ServiceDesc, srv)
}

func _RuleService_CreateRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServiceServer).CreateRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rulepb.RuleService/CreateRule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServiceServer).CreateRule(ctx, req.(*CreateRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleService_GetRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServiceServer).GetRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rulepb.RuleService/GetRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServiceServer).GetRules(ctx, req.(*GetRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleService_Classify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClassifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServiceServer).Classify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rulepb.RuleService/Classify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServiceServer).Classify(ctx, req.(*ClassifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RuleService_ServiceDesc is the grpc.ServiceDesc for RuleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RuleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rulepb.RuleService",
	HandlerType: (*RuleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRule",
			Handler:    _RuleService_CreateRule_Handler,
		},
		{
			MethodName: "GetRules",
			Handler:    _RuleService_GetRules_Handler,
		},
		{
			MethodName: "Classify",
			Handler:    _RuleService_Classify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rule.proto",
}
