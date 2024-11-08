// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: metainstance_service.proto

package mrdspb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MetaInstances_Create_FullMethodName                   = "/proto.mrds.ledger.metainstance.MetaInstances/Create"
	MetaInstances_GetByID_FullMethodName                  = "/proto.mrds.ledger.metainstance.MetaInstances/GetByID"
	MetaInstances_GetByName_FullMethodName                = "/proto.mrds.ledger.metainstance.MetaInstances/GetByName"
	MetaInstances_UpdateStatus_FullMethodName             = "/proto.mrds.ledger.metainstance.MetaInstances/UpdateStatus"
	MetaInstances_UpdateDeploymentID_FullMethodName       = "/proto.mrds.ledger.metainstance.MetaInstances/UpdateDeploymentID"
	MetaInstances_List_FullMethodName                     = "/proto.mrds.ledger.metainstance.MetaInstances/List"
	MetaInstances_Delete_FullMethodName                   = "/proto.mrds.ledger.metainstance.MetaInstances/Delete"
	MetaInstances_AddRuntimeInstance_FullMethodName       = "/proto.mrds.ledger.metainstance.MetaInstances/AddRuntimeInstance"
	MetaInstances_UpdateRuntimeStatus_FullMethodName      = "/proto.mrds.ledger.metainstance.MetaInstances/UpdateRuntimeStatus"
	MetaInstances_UpdateRuntimeActiveState_FullMethodName = "/proto.mrds.ledger.metainstance.MetaInstances/UpdateRuntimeActiveState"
	MetaInstances_RemoveRuntimeInstance_FullMethodName    = "/proto.mrds.ledger.metainstance.MetaInstances/RemoveRuntimeInstance"
	MetaInstances_AddOperation_FullMethodName             = "/proto.mrds.ledger.metainstance.MetaInstances/AddOperation"
	MetaInstances_UpdateOperationStatus_FullMethodName    = "/proto.mrds.ledger.metainstance.MetaInstances/UpdateOperationStatus"
	MetaInstances_RemoveOperation_FullMethodName          = "/proto.mrds.ledger.metainstance.MetaInstances/RemoveOperation"
)

// MetaInstancesClient is the client API for MetaInstances service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service definition for managing MetaInstance records.
type MetaInstancesClient interface {
	// Create a new MetaInstance.
	Create(ctx context.Context, in *CreateMetaInstanceRequest, opts ...grpc.CallOption) (*CreateMetaInstanceResponse, error)
	// Get a MetaInstance by its ID.
	GetByID(ctx context.Context, in *GetMetaInstanceByIDRequest, opts ...grpc.CallOption) (*GetMetaInstanceResponse, error)
	// Get a MetaInstance by its name.
	GetByName(ctx context.Context, in *GetMetaInstanceByNameRequest, opts ...grpc.CallOption) (*GetMetaInstanceResponse, error)
	// Update the status of an existing MetaInstance.
	UpdateStatus(ctx context.Context, in *UpdateMetaInstanceStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	// Update the DeploymentID of an existing MetaInstance.
	UpdateDeploymentID(ctx context.Context, in *UpdateDeploymentIDRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	// List MetaInstances that match the provided filters.
	List(ctx context.Context, in *ListMetaInstanceRequest, opts ...grpc.CallOption) (*ListMetaInstanceResponse, error)
	// Delete a MetaInstance by its metadata.
	Delete(ctx context.Context, in *DeleteMetaInstanceRequest, opts ...grpc.CallOption) (*DeleteMetaInstanceResponse, error)
	AddRuntimeInstance(ctx context.Context, in *AddRuntimeInstanceRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	UpdateRuntimeStatus(ctx context.Context, in *UpdateRuntimeStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	UpdateRuntimeActiveState(ctx context.Context, in *UpdateRuntimeActiveStateRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	RemoveRuntimeInstance(ctx context.Context, in *RemoveRuntimeInstanceRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	AddOperation(ctx context.Context, in *AddOperationRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	UpdateOperationStatus(ctx context.Context, in *UpdateOperationStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
	RemoveOperation(ctx context.Context, in *RemoveOperationRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error)
}

type metaInstancesClient struct {
	cc grpc.ClientConnInterface
}

func NewMetaInstancesClient(cc grpc.ClientConnInterface) MetaInstancesClient {
	return &metaInstancesClient{cc}
}

func (c *metaInstancesClient) Create(ctx context.Context, in *CreateMetaInstanceRequest, opts ...grpc.CallOption) (*CreateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) GetByID(ctx context.Context, in *GetMetaInstanceByIDRequest, opts ...grpc.CallOption) (*GetMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_GetByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) GetByName(ctx context.Context, in *GetMetaInstanceByNameRequest, opts ...grpc.CallOption) (*GetMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_GetByName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) UpdateStatus(ctx context.Context, in *UpdateMetaInstanceStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_UpdateStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) UpdateDeploymentID(ctx context.Context, in *UpdateDeploymentIDRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_UpdateDeploymentID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) List(ctx context.Context, in *ListMetaInstanceRequest, opts ...grpc.CallOption) (*ListMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) Delete(ctx context.Context, in *DeleteMetaInstanceRequest, opts ...grpc.CallOption) (*DeleteMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) AddRuntimeInstance(ctx context.Context, in *AddRuntimeInstanceRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_AddRuntimeInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) UpdateRuntimeStatus(ctx context.Context, in *UpdateRuntimeStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_UpdateRuntimeStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) UpdateRuntimeActiveState(ctx context.Context, in *UpdateRuntimeActiveStateRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_UpdateRuntimeActiveState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) RemoveRuntimeInstance(ctx context.Context, in *RemoveRuntimeInstanceRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_RemoveRuntimeInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) AddOperation(ctx context.Context, in *AddOperationRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_AddOperation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) UpdateOperationStatus(ctx context.Context, in *UpdateOperationStatusRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_UpdateOperationStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaInstancesClient) RemoveOperation(ctx context.Context, in *RemoveOperationRequest, opts ...grpc.CallOption) (*UpdateMetaInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetaInstanceResponse)
	err := c.cc.Invoke(ctx, MetaInstances_RemoveOperation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaInstancesServer is the server API for MetaInstances service.
// All implementations must embed UnimplementedMetaInstancesServer
// for forward compatibility.
//
// Service definition for managing MetaInstance records.
type MetaInstancesServer interface {
	// Create a new MetaInstance.
	Create(context.Context, *CreateMetaInstanceRequest) (*CreateMetaInstanceResponse, error)
	// Get a MetaInstance by its ID.
	GetByID(context.Context, *GetMetaInstanceByIDRequest) (*GetMetaInstanceResponse, error)
	// Get a MetaInstance by its name.
	GetByName(context.Context, *GetMetaInstanceByNameRequest) (*GetMetaInstanceResponse, error)
	// Update the status of an existing MetaInstance.
	UpdateStatus(context.Context, *UpdateMetaInstanceStatusRequest) (*UpdateMetaInstanceResponse, error)
	// Update the DeploymentID of an existing MetaInstance.
	UpdateDeploymentID(context.Context, *UpdateDeploymentIDRequest) (*UpdateMetaInstanceResponse, error)
	// List MetaInstances that match the provided filters.
	List(context.Context, *ListMetaInstanceRequest) (*ListMetaInstanceResponse, error)
	// Delete a MetaInstance by its metadata.
	Delete(context.Context, *DeleteMetaInstanceRequest) (*DeleteMetaInstanceResponse, error)
	AddRuntimeInstance(context.Context, *AddRuntimeInstanceRequest) (*UpdateMetaInstanceResponse, error)
	UpdateRuntimeStatus(context.Context, *UpdateRuntimeStatusRequest) (*UpdateMetaInstanceResponse, error)
	UpdateRuntimeActiveState(context.Context, *UpdateRuntimeActiveStateRequest) (*UpdateMetaInstanceResponse, error)
	RemoveRuntimeInstance(context.Context, *RemoveRuntimeInstanceRequest) (*UpdateMetaInstanceResponse, error)
	AddOperation(context.Context, *AddOperationRequest) (*UpdateMetaInstanceResponse, error)
	UpdateOperationStatus(context.Context, *UpdateOperationStatusRequest) (*UpdateMetaInstanceResponse, error)
	RemoveOperation(context.Context, *RemoveOperationRequest) (*UpdateMetaInstanceResponse, error)
	mustEmbedUnimplementedMetaInstancesServer()
}

// UnimplementedMetaInstancesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetaInstancesServer struct{}

func (UnimplementedMetaInstancesServer) Create(context.Context, *CreateMetaInstanceRequest) (*CreateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedMetaInstancesServer) GetByID(context.Context, *GetMetaInstanceByIDRequest) (*GetMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedMetaInstancesServer) GetByName(context.Context, *GetMetaInstanceByNameRequest) (*GetMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedMetaInstancesServer) UpdateStatus(context.Context, *UpdateMetaInstanceStatusRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedMetaInstancesServer) UpdateDeploymentID(context.Context, *UpdateDeploymentIDRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDeploymentID not implemented")
}
func (UnimplementedMetaInstancesServer) List(context.Context, *ListMetaInstanceRequest) (*ListMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedMetaInstancesServer) Delete(context.Context, *DeleteMetaInstanceRequest) (*DeleteMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedMetaInstancesServer) AddRuntimeInstance(context.Context, *AddRuntimeInstanceRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRuntimeInstance not implemented")
}
func (UnimplementedMetaInstancesServer) UpdateRuntimeStatus(context.Context, *UpdateRuntimeStatusRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRuntimeStatus not implemented")
}
func (UnimplementedMetaInstancesServer) UpdateRuntimeActiveState(context.Context, *UpdateRuntimeActiveStateRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRuntimeActiveState not implemented")
}
func (UnimplementedMetaInstancesServer) RemoveRuntimeInstance(context.Context, *RemoveRuntimeInstanceRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveRuntimeInstance not implemented")
}
func (UnimplementedMetaInstancesServer) AddOperation(context.Context, *AddOperationRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOperation not implemented")
}
func (UnimplementedMetaInstancesServer) UpdateOperationStatus(context.Context, *UpdateOperationStatusRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOperationStatus not implemented")
}
func (UnimplementedMetaInstancesServer) RemoveOperation(context.Context, *RemoveOperationRequest) (*UpdateMetaInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOperation not implemented")
}
func (UnimplementedMetaInstancesServer) mustEmbedUnimplementedMetaInstancesServer() {}
func (UnimplementedMetaInstancesServer) testEmbeddedByValue()                       {}

// UnsafeMetaInstancesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetaInstancesServer will
// result in compilation errors.
type UnsafeMetaInstancesServer interface {
	mustEmbedUnimplementedMetaInstancesServer()
}

func RegisterMetaInstancesServer(s grpc.ServiceRegistrar, srv MetaInstancesServer) {
	// If the following call pancis, it indicates UnimplementedMetaInstancesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MetaInstances_ServiceDesc, srv)
}

func _MetaInstances_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMetaInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).Create(ctx, req.(*CreateMetaInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetaInstanceByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_GetByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).GetByID(ctx, req.(*GetMetaInstanceByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetaInstanceByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_GetByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).GetByName(ctx, req.(*GetMetaInstanceByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetaInstanceStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_UpdateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).UpdateStatus(ctx, req.(*UpdateMetaInstanceStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_UpdateDeploymentID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDeploymentIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).UpdateDeploymentID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_UpdateDeploymentID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).UpdateDeploymentID(ctx, req.(*UpdateDeploymentIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMetaInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).List(ctx, req.(*ListMetaInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMetaInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).Delete(ctx, req.(*DeleteMetaInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_AddRuntimeInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRuntimeInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).AddRuntimeInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_AddRuntimeInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).AddRuntimeInstance(ctx, req.(*AddRuntimeInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_UpdateRuntimeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRuntimeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).UpdateRuntimeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_UpdateRuntimeStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).UpdateRuntimeStatus(ctx, req.(*UpdateRuntimeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_UpdateRuntimeActiveState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRuntimeActiveStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).UpdateRuntimeActiveState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_UpdateRuntimeActiveState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).UpdateRuntimeActiveState(ctx, req.(*UpdateRuntimeActiveStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_RemoveRuntimeInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRuntimeInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).RemoveRuntimeInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_RemoveRuntimeInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).RemoveRuntimeInstance(ctx, req.(*RemoveRuntimeInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_AddOperation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddOperationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).AddOperation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_AddOperation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).AddOperation(ctx, req.(*AddOperationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_UpdateOperationStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOperationStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).UpdateOperationStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_UpdateOperationStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).UpdateOperationStatus(ctx, req.(*UpdateOperationStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetaInstances_RemoveOperation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveOperationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaInstancesServer).RemoveOperation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetaInstances_RemoveOperation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaInstancesServer).RemoveOperation(ctx, req.(*RemoveOperationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetaInstances_ServiceDesc is the grpc.ServiceDesc for MetaInstances service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetaInstances_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.mrds.ledger.metainstance.MetaInstances",
	HandlerType: (*MetaInstancesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _MetaInstances_Create_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _MetaInstances_GetByID_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _MetaInstances_GetByName_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _MetaInstances_UpdateStatus_Handler,
		},
		{
			MethodName: "UpdateDeploymentID",
			Handler:    _MetaInstances_UpdateDeploymentID_Handler,
		},
		{
			MethodName: "List",
			Handler:    _MetaInstances_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _MetaInstances_Delete_Handler,
		},
		{
			MethodName: "AddRuntimeInstance",
			Handler:    _MetaInstances_AddRuntimeInstance_Handler,
		},
		{
			MethodName: "UpdateRuntimeStatus",
			Handler:    _MetaInstances_UpdateRuntimeStatus_Handler,
		},
		{
			MethodName: "UpdateRuntimeActiveState",
			Handler:    _MetaInstances_UpdateRuntimeActiveState_Handler,
		},
		{
			MethodName: "RemoveRuntimeInstance",
			Handler:    _MetaInstances_RemoveRuntimeInstance_Handler,
		},
		{
			MethodName: "AddOperation",
			Handler:    _MetaInstances_AddOperation_Handler,
		},
		{
			MethodName: "UpdateOperationStatus",
			Handler:    _MetaInstances_UpdateOperationStatus_Handler,
		},
		{
			MethodName: "RemoveOperation",
			Handler:    _MetaInstances_RemoveOperation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metainstance_service.proto",
}
