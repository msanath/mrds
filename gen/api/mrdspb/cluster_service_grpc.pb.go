// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: cluster_service.proto

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
	Clusters_Create_FullMethodName       = "/proto.mrds.ledger.cluster.Clusters/Create"
	Clusters_GetByID_FullMethodName      = "/proto.mrds.ledger.cluster.Clusters/GetByID"
	Clusters_GetByName_FullMethodName    = "/proto.mrds.ledger.cluster.Clusters/GetByName"
	Clusters_UpdateStatus_FullMethodName = "/proto.mrds.ledger.cluster.Clusters/UpdateStatus"
	Clusters_List_FullMethodName         = "/proto.mrds.ledger.cluster.Clusters/List"
	Clusters_Delete_FullMethodName       = "/proto.mrds.ledger.cluster.Clusters/Delete"
)

// ClustersClient is the client API for Clusters service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service definition for managing Cluster records.
type ClustersClient interface {
	// Create a new Cluster.
	Create(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*CreateClusterResponse, error)
	// Get a Cluster by its ID.
	GetByID(ctx context.Context, in *GetClusterByIDRequest, opts ...grpc.CallOption) (*GetClusterResponse, error)
	// Get a Cluster by its name.
	GetByName(ctx context.Context, in *GetClusterByNameRequest, opts ...grpc.CallOption) (*GetClusterResponse, error)
	// Update the state of an existing Cluster.
	UpdateStatus(ctx context.Context, in *UpdateClusterStatusRequest, opts ...grpc.CallOption) (*UpdateClusterResponse, error)
	// List Clusters that match the provided filters.
	List(ctx context.Context, in *ListClusterRequest, opts ...grpc.CallOption) (*ListClusterResponse, error)
	// Delete a Cluster by its metadata.
	Delete(ctx context.Context, in *DeleteClusterRequest, opts ...grpc.CallOption) (*DeleteClusterResponse, error)
}

type clustersClient struct {
	cc grpc.ClientConnInterface
}

func NewClustersClient(cc grpc.ClientConnInterface) ClustersClient {
	return &clustersClient{cc}
}

func (c *clustersClient) Create(ctx context.Context, in *CreateClusterRequest, opts ...grpc.CallOption) (*CreateClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clustersClient) GetByID(ctx context.Context, in *GetClusterByIDRequest, opts ...grpc.CallOption) (*GetClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_GetByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clustersClient) GetByName(ctx context.Context, in *GetClusterByNameRequest, opts ...grpc.CallOption) (*GetClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_GetByName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clustersClient) UpdateStatus(ctx context.Context, in *UpdateClusterStatusRequest, opts ...grpc.CallOption) (*UpdateClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_UpdateStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clustersClient) List(ctx context.Context, in *ListClusterRequest, opts ...grpc.CallOption) (*ListClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clustersClient) Delete(ctx context.Context, in *DeleteClusterRequest, opts ...grpc.CallOption) (*DeleteClusterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteClusterResponse)
	err := c.cc.Invoke(ctx, Clusters_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClustersServer is the server API for Clusters service.
// All implementations must embed UnimplementedClustersServer
// for forward compatibility.
//
// Service definition for managing Cluster records.
type ClustersServer interface {
	// Create a new Cluster.
	Create(context.Context, *CreateClusterRequest) (*CreateClusterResponse, error)
	// Get a Cluster by its ID.
	GetByID(context.Context, *GetClusterByIDRequest) (*GetClusterResponse, error)
	// Get a Cluster by its name.
	GetByName(context.Context, *GetClusterByNameRequest) (*GetClusterResponse, error)
	// Update the state of an existing Cluster.
	UpdateStatus(context.Context, *UpdateClusterStatusRequest) (*UpdateClusterResponse, error)
	// List Clusters that match the provided filters.
	List(context.Context, *ListClusterRequest) (*ListClusterResponse, error)
	// Delete a Cluster by its metadata.
	Delete(context.Context, *DeleteClusterRequest) (*DeleteClusterResponse, error)
	mustEmbedUnimplementedClustersServer()
}

// UnimplementedClustersServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClustersServer struct{}

func (UnimplementedClustersServer) Create(context.Context, *CreateClusterRequest) (*CreateClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedClustersServer) GetByID(context.Context, *GetClusterByIDRequest) (*GetClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
func (UnimplementedClustersServer) GetByName(context.Context, *GetClusterByNameRequest) (*GetClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedClustersServer) UpdateStatus(context.Context, *UpdateClusterStatusRequest) (*UpdateClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedClustersServer) List(context.Context, *ListClusterRequest) (*ListClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedClustersServer) Delete(context.Context, *DeleteClusterRequest) (*DeleteClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedClustersServer) mustEmbedUnimplementedClustersServer() {}
func (UnimplementedClustersServer) testEmbeddedByValue()                  {}

// UnsafeClustersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClustersServer will
// result in compilation errors.
type UnsafeClustersServer interface {
	mustEmbedUnimplementedClustersServer()
}

func RegisterClustersServer(s grpc.ServiceRegistrar, srv ClustersServer) {
	// If the following call pancis, it indicates UnimplementedClustersServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Clusters_ServiceDesc, srv)
}

func _Clusters_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).Create(ctx, req.(*CreateClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Clusters_GetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).GetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_GetByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).GetByID(ctx, req.(*GetClusterByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Clusters_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_GetByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).GetByName(ctx, req.(*GetClusterByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Clusters_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateClusterStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_UpdateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).UpdateStatus(ctx, req.(*UpdateClusterStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Clusters_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).List(ctx, req.(*ListClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Clusters_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClustersServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Clusters_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClustersServer).Delete(ctx, req.(*DeleteClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Clusters_ServiceDesc is the grpc.ServiceDesc for Clusters service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Clusters_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.mrds.ledger.cluster.Clusters",
	HandlerType: (*ClustersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Clusters_Create_Handler,
		},
		{
			MethodName: "GetByID",
			Handler:    _Clusters_GetByID_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _Clusters_GetByName_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _Clusters_UpdateStatus_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Clusters_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Clusters_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster_service.proto",
}
