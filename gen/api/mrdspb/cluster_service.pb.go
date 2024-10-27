// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: cluster_service.proto

package mrdspb

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

// Request to create a new Cluster.
type CreateClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *CreateClusterRequest) Reset() {
	*x = CreateClusterRequest{}
	mi := &file_cluster_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClusterRequest) ProtoMessage() {}

func (x *CreateClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClusterRequest.ProtoReflect.Descriptor instead.
func (*CreateClusterRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateClusterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response after creating a new Cluster.
type CreateClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The newly created Cluster record.
	Record *Cluster `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *CreateClusterResponse) Reset() {
	*x = CreateClusterResponse{}
	mi := &file_cluster_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClusterResponse) ProtoMessage() {}

func (x *CreateClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClusterResponse.ProtoReflect.Descriptor instead.
func (*CreateClusterResponse) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateClusterResponse) GetRecord() *Cluster {
	if x != nil {
		return x.Record
	}
	return nil
}

// Request to update the state and message of a Cluster.
type UpdateClusterStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metadata of the Cluster to update.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The new state of the Cluster.
	Status *ClusterStatus `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *UpdateClusterStatusRequest) Reset() {
	*x = UpdateClusterStatusRequest{}
	mi := &file_cluster_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateClusterStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateClusterStatusRequest) ProtoMessage() {}

func (x *UpdateClusterStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateClusterStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateClusterStatusRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateClusterStatusRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *UpdateClusterStatusRequest) GetStatus() *ClusterStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// Response after updating the state of a Cluster.
type UpdateClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The updated Cluster record.
	Record *Cluster `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *UpdateClusterResponse) Reset() {
	*x = UpdateClusterResponse{}
	mi := &file_cluster_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateClusterResponse) ProtoMessage() {}

func (x *UpdateClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateClusterResponse.ProtoReflect.Descriptor instead.
func (*UpdateClusterResponse) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateClusterResponse) GetRecord() *Cluster {
	if x != nil {
		return x.Record
	}
	return nil
}

type GetClusterByMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *GetClusterByMetadataRequest) Reset() {
	*x = GetClusterByMetadataRequest{}
	mi := &file_cluster_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClusterByMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterByMetadataRequest) ProtoMessage() {}

func (x *GetClusterByMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterByMetadataRequest.ProtoReflect.Descriptor instead.
func (*GetClusterByMetadataRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetClusterByMetadataRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Request for getting a Cluster by its name.
type GetClusterByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the Cluster to get.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetClusterByNameRequest) Reset() {
	*x = GetClusterByNameRequest{}
	mi := &file_cluster_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClusterByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterByNameRequest) ProtoMessage() {}

func (x *GetClusterByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterByNameRequest.ProtoReflect.Descriptor instead.
func (*GetClusterByNameRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{5}
}

func (x *GetClusterByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response after fetching a Cluster.
type GetClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Cluster record that was fetched.
	Record *Cluster `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *GetClusterResponse) Reset() {
	*x = GetClusterResponse{}
	mi := &file_cluster_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClusterResponse) ProtoMessage() {}

func (x *GetClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClusterResponse.ProtoReflect.Descriptor instead.
func (*GetClusterResponse) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{6}
}

func (x *GetClusterResponse) GetRecord() *Cluster {
	if x != nil {
		return x.Record
	}
	return nil
}

// Request to list Clusters with specific filters.
type ListClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// IN condition for filtering by IDs.
	IdIn []string `protobuf:"bytes,1,rep,name=id_in,json=idIn,proto3" json:"id_in,omitempty"`
	// IN condition for filtering by Names.
	NameIn []string `protobuf:"bytes,2,rep,name=name_in,json=nameIn,proto3" json:"name_in,omitempty"`
	// Greater than or equal condition for filtering by version.
	VersionGte uint64 `protobuf:"varint,3,opt,name=version_gte,json=versionGte,proto3" json:"version_gte,omitempty"`
	// Less than or equal condition for filtering by version.
	VersionLte uint64 `protobuf:"varint,4,opt,name=version_lte,json=versionLte,proto3" json:"version_lte,omitempty"`
	// Equal condition for filtering by version.
	VersionEq uint64 `protobuf:"varint,5,opt,name=version_eq,json=versionEq,proto3" json:"version_eq,omitempty"`
	// IN condition for filtering by state.
	StateIn []ClusterState `protobuf:"varint,6,rep,packed,name=state_in,json=stateIn,proto3,enum=proto.mrds.ledger.cluster.ClusterState" json:"state_in,omitempty"`
	// NOT IN condition for filtering by state.
	StateNotIn []ClusterState `protobuf:"varint,7,rep,packed,name=state_not_in,json=stateNotIn,proto3,enum=proto.mrds.ledger.cluster.ClusterState" json:"state_not_in,omitempty"`
	// Whether to include soft-deleted resources in the query.
	IncludeDeleted bool `protobuf:"varint,8,opt,name=include_deleted,json=includeDeleted,proto3" json:"include_deleted,omitempty"`
	// Limit the number of results returned.
	Limit uint32 `protobuf:"varint,9,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListClusterRequest) Reset() {
	*x = ListClusterRequest{}
	mi := &file_cluster_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClusterRequest) ProtoMessage() {}

func (x *ListClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClusterRequest.ProtoReflect.Descriptor instead.
func (*ListClusterRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{7}
}

func (x *ListClusterRequest) GetIdIn() []string {
	if x != nil {
		return x.IdIn
	}
	return nil
}

func (x *ListClusterRequest) GetNameIn() []string {
	if x != nil {
		return x.NameIn
	}
	return nil
}

func (x *ListClusterRequest) GetVersionGte() uint64 {
	if x != nil {
		return x.VersionGte
	}
	return 0
}

func (x *ListClusterRequest) GetVersionLte() uint64 {
	if x != nil {
		return x.VersionLte
	}
	return 0
}

func (x *ListClusterRequest) GetVersionEq() uint64 {
	if x != nil {
		return x.VersionEq
	}
	return 0
}

func (x *ListClusterRequest) GetStateIn() []ClusterState {
	if x != nil {
		return x.StateIn
	}
	return nil
}

func (x *ListClusterRequest) GetStateNotIn() []ClusterState {
	if x != nil {
		return x.StateNotIn
	}
	return nil
}

func (x *ListClusterRequest) GetIncludeDeleted() bool {
	if x != nil {
		return x.IncludeDeleted
	}
	return false
}

func (x *ListClusterRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// Response for listing Clusters.
type ListClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of Cluster records that match the query.
	Records []*Cluster `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *ListClusterResponse) Reset() {
	*x = ListClusterResponse{}
	mi := &file_cluster_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClusterResponse) ProtoMessage() {}

func (x *ListClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClusterResponse.ProtoReflect.Descriptor instead.
func (*ListClusterResponse) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{8}
}

func (x *ListClusterResponse) GetRecords() []*Cluster {
	if x != nil {
		return x.Records
	}
	return nil
}

// Request to delete a Cluster by its metadata.
type DeleteClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metadata of the Cluster to delete.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *DeleteClusterRequest) Reset() {
	*x = DeleteClusterRequest{}
	mi := &file_cluster_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteClusterRequest) ProtoMessage() {}

func (x *DeleteClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteClusterRequest.ProtoReflect.Descriptor instead.
func (*DeleteClusterRequest) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteClusterRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Response after deleting a Cluster.
type DeleteClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteClusterResponse) Reset() {
	*x = DeleteClusterResponse{}
	mi := &file_cluster_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteClusterResponse) ProtoMessage() {}

func (x *DeleteClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cluster_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteClusterResponse.ProtoReflect.Descriptor instead.
func (*DeleteClusterResponse) Descriptor() ([]byte, []int) {
	return file_cluster_service_proto_rawDescGZIP(), []int{10}
}

var File_cluster_service_proto protoreflect.FileDescriptor

var file_cluster_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x1a, 0x0e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0d, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2a, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x53, 0x0a,
	0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x22, 0x95, 0x01, 0x0a, 0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x40, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x53, 0x0a, 0x15, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73,
	0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22,
	0x54, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x79, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2d, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x50, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x06, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x06,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0xf1, 0x02, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x13, 0x0a,
	0x05, 0x69, 0x64, 0x5f, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x69, 0x64,
	0x49, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x67, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x47, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4c, 0x74, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x71, 0x12, 0x42, 0x0a, 0x08,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x27,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x07, 0x73, 0x74, 0x61, 0x74, 0x65, 0x49, 0x6e,
	0x12, 0x49, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x6e, 0x6f, 0x74, 0x5f, 0x69, 0x6e,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x69,
	0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x53, 0x0a, 0x13, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3c, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e,
	0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22,
	0x4d, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x17,
	0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xac, 0x05, 0x0a, 0x08, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x73, 0x12, 0x6b, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x2f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x76, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x36, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e,
	0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6e, 0x0a, 0x09, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x42, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x77, 0x0a, 0x0c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x35, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x65, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6b, 0x0a, 0x06, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x12, 0x2f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73,
	0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64,
	0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d,
	0x72, 0x64, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cluster_service_proto_rawDescOnce sync.Once
	file_cluster_service_proto_rawDescData = file_cluster_service_proto_rawDesc
)

func file_cluster_service_proto_rawDescGZIP() []byte {
	file_cluster_service_proto_rawDescOnce.Do(func() {
		file_cluster_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_cluster_service_proto_rawDescData)
	})
	return file_cluster_service_proto_rawDescData
}

var file_cluster_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_cluster_service_proto_goTypes = []any{
	(*CreateClusterRequest)(nil),        // 0: proto.mrds.ledger.cluster.CreateClusterRequest
	(*CreateClusterResponse)(nil),       // 1: proto.mrds.ledger.cluster.CreateClusterResponse
	(*UpdateClusterStatusRequest)(nil),  // 2: proto.mrds.ledger.cluster.UpdateClusterStatusRequest
	(*UpdateClusterResponse)(nil),       // 3: proto.mrds.ledger.cluster.UpdateClusterResponse
	(*GetClusterByMetadataRequest)(nil), // 4: proto.mrds.ledger.cluster.GetClusterByMetadataRequest
	(*GetClusterByNameRequest)(nil),     // 5: proto.mrds.ledger.cluster.GetClusterByNameRequest
	(*GetClusterResponse)(nil),          // 6: proto.mrds.ledger.cluster.GetClusterResponse
	(*ListClusterRequest)(nil),          // 7: proto.mrds.ledger.cluster.ListClusterRequest
	(*ListClusterResponse)(nil),         // 8: proto.mrds.ledger.cluster.ListClusterResponse
	(*DeleteClusterRequest)(nil),        // 9: proto.mrds.ledger.cluster.DeleteClusterRequest
	(*DeleteClusterResponse)(nil),       // 10: proto.mrds.ledger.cluster.DeleteClusterResponse
	(*Cluster)(nil),                     // 11: proto.mrds.ledger.cluster.Cluster
	(*Metadata)(nil),                    // 12: proto.mrds.core.Metadata
	(*ClusterStatus)(nil),               // 13: proto.mrds.ledger.cluster.ClusterStatus
	(ClusterState)(0),                   // 14: proto.mrds.ledger.cluster.ClusterState
}
var file_cluster_service_proto_depIdxs = []int32{
	11, // 0: proto.mrds.ledger.cluster.CreateClusterResponse.record:type_name -> proto.mrds.ledger.cluster.Cluster
	12, // 1: proto.mrds.ledger.cluster.UpdateClusterStatusRequest.metadata:type_name -> proto.mrds.core.Metadata
	13, // 2: proto.mrds.ledger.cluster.UpdateClusterStatusRequest.status:type_name -> proto.mrds.ledger.cluster.ClusterStatus
	11, // 3: proto.mrds.ledger.cluster.UpdateClusterResponse.record:type_name -> proto.mrds.ledger.cluster.Cluster
	12, // 4: proto.mrds.ledger.cluster.GetClusterByMetadataRequest.metadata:type_name -> proto.mrds.core.Metadata
	11, // 5: proto.mrds.ledger.cluster.GetClusterResponse.record:type_name -> proto.mrds.ledger.cluster.Cluster
	14, // 6: proto.mrds.ledger.cluster.ListClusterRequest.state_in:type_name -> proto.mrds.ledger.cluster.ClusterState
	14, // 7: proto.mrds.ledger.cluster.ListClusterRequest.state_not_in:type_name -> proto.mrds.ledger.cluster.ClusterState
	11, // 8: proto.mrds.ledger.cluster.ListClusterResponse.records:type_name -> proto.mrds.ledger.cluster.Cluster
	12, // 9: proto.mrds.ledger.cluster.DeleteClusterRequest.metadata:type_name -> proto.mrds.core.Metadata
	0,  // 10: proto.mrds.ledger.cluster.Clusters.Create:input_type -> proto.mrds.ledger.cluster.CreateClusterRequest
	4,  // 11: proto.mrds.ledger.cluster.Clusters.GetByMetadata:input_type -> proto.mrds.ledger.cluster.GetClusterByMetadataRequest
	5,  // 12: proto.mrds.ledger.cluster.Clusters.GetByName:input_type -> proto.mrds.ledger.cluster.GetClusterByNameRequest
	2,  // 13: proto.mrds.ledger.cluster.Clusters.UpdateStatus:input_type -> proto.mrds.ledger.cluster.UpdateClusterStatusRequest
	7,  // 14: proto.mrds.ledger.cluster.Clusters.List:input_type -> proto.mrds.ledger.cluster.ListClusterRequest
	9,  // 15: proto.mrds.ledger.cluster.Clusters.Delete:input_type -> proto.mrds.ledger.cluster.DeleteClusterRequest
	1,  // 16: proto.mrds.ledger.cluster.Clusters.Create:output_type -> proto.mrds.ledger.cluster.CreateClusterResponse
	6,  // 17: proto.mrds.ledger.cluster.Clusters.GetByMetadata:output_type -> proto.mrds.ledger.cluster.GetClusterResponse
	6,  // 18: proto.mrds.ledger.cluster.Clusters.GetByName:output_type -> proto.mrds.ledger.cluster.GetClusterResponse
	3,  // 19: proto.mrds.ledger.cluster.Clusters.UpdateStatus:output_type -> proto.mrds.ledger.cluster.UpdateClusterResponse
	8,  // 20: proto.mrds.ledger.cluster.Clusters.List:output_type -> proto.mrds.ledger.cluster.ListClusterResponse
	10, // 21: proto.mrds.ledger.cluster.Clusters.Delete:output_type -> proto.mrds.ledger.cluster.DeleteClusterResponse
	16, // [16:22] is the sub-list for method output_type
	10, // [10:16] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_cluster_service_proto_init() }
func file_cluster_service_proto_init() {
	if File_cluster_service_proto != nil {
		return
	}
	file_metadata_proto_init()
	file_cluster_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cluster_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cluster_service_proto_goTypes,
		DependencyIndexes: file_cluster_service_proto_depIdxs,
		MessageInfos:      file_cluster_service_proto_msgTypes,
	}.Build()
	File_cluster_service_proto = out.File
	file_cluster_service_proto_rawDesc = nil
	file_cluster_service_proto_goTypes = nil
	file_cluster_service_proto_depIdxs = nil
}
