// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: metainstance_service.proto

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

// Request to create a new MetaInstance.
type CreateMetaInstanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *CreateMetaInstanceRequest) Reset() {
	*x = CreateMetaInstanceRequest{}
	mi := &file_metainstance_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMetaInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMetaInstanceRequest) ProtoMessage() {}

func (x *CreateMetaInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMetaInstanceRequest.ProtoReflect.Descriptor instead.
func (*CreateMetaInstanceRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateMetaInstanceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response after creating a new MetaInstance.
type CreateMetaInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The newly created MetaInstance record.
	Record *MetaInstance `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *CreateMetaInstanceResponse) Reset() {
	*x = CreateMetaInstanceResponse{}
	mi := &file_metainstance_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMetaInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMetaInstanceResponse) ProtoMessage() {}

func (x *CreateMetaInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMetaInstanceResponse.ProtoReflect.Descriptor instead.
func (*CreateMetaInstanceResponse) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateMetaInstanceResponse) GetRecord() *MetaInstance {
	if x != nil {
		return x.Record
	}
	return nil
}

// Request to update the state and message of a MetaInstance.
type UpdateMetaInstanceStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metadata of the MetaInstance to update.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The new state of the MetaInstance.
	Status *MetaInstanceStatus `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *UpdateMetaInstanceStatusRequest) Reset() {
	*x = UpdateMetaInstanceStatusRequest{}
	mi := &file_metainstance_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateMetaInstanceStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetaInstanceStatusRequest) ProtoMessage() {}

func (x *UpdateMetaInstanceStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetaInstanceStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetaInstanceStatusRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateMetaInstanceStatusRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *UpdateMetaInstanceStatusRequest) GetStatus() *MetaInstanceStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// Response after updating the state of a MetaInstance.
type UpdateMetaInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The updated MetaInstance record.
	Record *MetaInstance `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *UpdateMetaInstanceResponse) Reset() {
	*x = UpdateMetaInstanceResponse{}
	mi := &file_metainstance_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateMetaInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetaInstanceResponse) ProtoMessage() {}

func (x *UpdateMetaInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetaInstanceResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetaInstanceResponse) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateMetaInstanceResponse) GetRecord() *MetaInstance {
	if x != nil {
		return x.Record
	}
	return nil
}

type GetMetaInstanceByMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *GetMetaInstanceByMetadataRequest) Reset() {
	*x = GetMetaInstanceByMetadataRequest{}
	mi := &file_metainstance_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetaInstanceByMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaInstanceByMetadataRequest) ProtoMessage() {}

func (x *GetMetaInstanceByMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaInstanceByMetadataRequest.ProtoReflect.Descriptor instead.
func (*GetMetaInstanceByMetadataRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetMetaInstanceByMetadataRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Request for getting a MetaInstance by its name.
type GetMetaInstanceByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the MetaInstance to get.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetMetaInstanceByNameRequest) Reset() {
	*x = GetMetaInstanceByNameRequest{}
	mi := &file_metainstance_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetaInstanceByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaInstanceByNameRequest) ProtoMessage() {}

func (x *GetMetaInstanceByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaInstanceByNameRequest.ProtoReflect.Descriptor instead.
func (*GetMetaInstanceByNameRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{5}
}

func (x *GetMetaInstanceByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response after fetching a MetaInstance.
type GetMetaInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The MetaInstance record that was fetched.
	Record *MetaInstance `protobuf:"bytes,1,opt,name=record,proto3" json:"record,omitempty"`
}

func (x *GetMetaInstanceResponse) Reset() {
	*x = GetMetaInstanceResponse{}
	mi := &file_metainstance_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetaInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaInstanceResponse) ProtoMessage() {}

func (x *GetMetaInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaInstanceResponse.ProtoReflect.Descriptor instead.
func (*GetMetaInstanceResponse) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{6}
}

func (x *GetMetaInstanceResponse) GetRecord() *MetaInstance {
	if x != nil {
		return x.Record
	}
	return nil
}

// Request to list MetaInstances with specific filters.
type ListMetaInstanceRequest struct {
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
	StateIn []MetaInstanceState `protobuf:"varint,6,rep,packed,name=state_in,json=stateIn,proto3,enum=proto.mrds.ledger.metainstance.MetaInstanceState" json:"state_in,omitempty"`
	// NOT IN condition for filtering by state.
	StateNotIn []MetaInstanceState `protobuf:"varint,7,rep,packed,name=state_not_in,json=stateNotIn,proto3,enum=proto.mrds.ledger.metainstance.MetaInstanceState" json:"state_not_in,omitempty"`
	// Whether to include soft-deleted resources in the query.
	IncludeDeleted bool `protobuf:"varint,8,opt,name=include_deleted,json=includeDeleted,proto3" json:"include_deleted,omitempty"`
	// Limit the number of results returned.
	Limit uint32 `protobuf:"varint,9,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListMetaInstanceRequest) Reset() {
	*x = ListMetaInstanceRequest{}
	mi := &file_metainstance_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMetaInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetaInstanceRequest) ProtoMessage() {}

func (x *ListMetaInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetaInstanceRequest.ProtoReflect.Descriptor instead.
func (*ListMetaInstanceRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{7}
}

func (x *ListMetaInstanceRequest) GetIdIn() []string {
	if x != nil {
		return x.IdIn
	}
	return nil
}

func (x *ListMetaInstanceRequest) GetNameIn() []string {
	if x != nil {
		return x.NameIn
	}
	return nil
}

func (x *ListMetaInstanceRequest) GetVersionGte() uint64 {
	if x != nil {
		return x.VersionGte
	}
	return 0
}

func (x *ListMetaInstanceRequest) GetVersionLte() uint64 {
	if x != nil {
		return x.VersionLte
	}
	return 0
}

func (x *ListMetaInstanceRequest) GetVersionEq() uint64 {
	if x != nil {
		return x.VersionEq
	}
	return 0
}

func (x *ListMetaInstanceRequest) GetStateIn() []MetaInstanceState {
	if x != nil {
		return x.StateIn
	}
	return nil
}

func (x *ListMetaInstanceRequest) GetStateNotIn() []MetaInstanceState {
	if x != nil {
		return x.StateNotIn
	}
	return nil
}

func (x *ListMetaInstanceRequest) GetIncludeDeleted() bool {
	if x != nil {
		return x.IncludeDeleted
	}
	return false
}

func (x *ListMetaInstanceRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// Response for listing MetaInstances.
type ListMetaInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of MetaInstance records that match the query.
	Records []*MetaInstance `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *ListMetaInstanceResponse) Reset() {
	*x = ListMetaInstanceResponse{}
	mi := &file_metainstance_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMetaInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMetaInstanceResponse) ProtoMessage() {}

func (x *ListMetaInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMetaInstanceResponse.ProtoReflect.Descriptor instead.
func (*ListMetaInstanceResponse) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{8}
}

func (x *ListMetaInstanceResponse) GetRecords() []*MetaInstance {
	if x != nil {
		return x.Records
	}
	return nil
}

// Request to delete a MetaInstance by its metadata.
type DeleteMetaInstanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The metadata of the MetaInstance to delete.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *DeleteMetaInstanceRequest) Reset() {
	*x = DeleteMetaInstanceRequest{}
	mi := &file_metainstance_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteMetaInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMetaInstanceRequest) ProtoMessage() {}

func (x *DeleteMetaInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMetaInstanceRequest.ProtoReflect.Descriptor instead.
func (*DeleteMetaInstanceRequest) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteMetaInstanceRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Response after deleting a MetaInstance.
type DeleteMetaInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteMetaInstanceResponse) Reset() {
	*x = DeleteMetaInstanceResponse{}
	mi := &file_metainstance_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteMetaInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMetaInstanceResponse) ProtoMessage() {}

func (x *DeleteMetaInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMetaInstanceResponse.ProtoReflect.Descriptor instead.
func (*DeleteMetaInstanceResponse) Descriptor() ([]byte, []int) {
	return file_metainstance_service_proto_rawDescGZIP(), []int{10}
}

var File_metainstance_service_proto protoreflect.FileDescriptor

var file_metainstance_service_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e,
	0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x1a, 0x0e, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x6d, 0x65,
	0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2f, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x62, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x44, 0x0a, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0xa4, 0x01, 0x0a, 0x1f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x4a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x62, 0x0a, 0x1a,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x06, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x22, 0x59, 0x0a, 0x20, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x42, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x32, 0x0a, 0x1c, 0x47,
	0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x42, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x5f, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x06, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x22, 0x8a, 0x03, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x13, 0x0a, 0x05,
	0x69, 0x64, 0x5f, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x69, 0x64, 0x49,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x6e, 0x61, 0x6d, 0x65, 0x49, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x67, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x47, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0a, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4c, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x71, 0x12, 0x4c, 0x0a, 0x08, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x31, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x07, 0x73, 0x74, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x12, 0x53, 0x0a, 0x0c, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x5f, 0x6e, 0x6f, 0x74, 0x5f, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0e, 0x32,
	0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x12, 0x27,
	0x0a, 0x0f, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x62, 0x0a,
	0x18, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x07, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x73, 0x22, 0x52, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x1c, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xac, 0x06, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x7f, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x39, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8a, 0x01, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x42, 0x79,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x40, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74,
	0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74,
	0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x42, 0x79, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x3c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x37, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8b, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65,
	0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x79, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x37,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x38, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x74,
	0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x7f, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x39, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e,
	0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x72, 0x64, 0x73, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metainstance_service_proto_rawDescOnce sync.Once
	file_metainstance_service_proto_rawDescData = file_metainstance_service_proto_rawDesc
)

func file_metainstance_service_proto_rawDescGZIP() []byte {
	file_metainstance_service_proto_rawDescOnce.Do(func() {
		file_metainstance_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_metainstance_service_proto_rawDescData)
	})
	return file_metainstance_service_proto_rawDescData
}

var file_metainstance_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_metainstance_service_proto_goTypes = []any{
	(*CreateMetaInstanceRequest)(nil),        // 0: proto.mrds.ledger.metainstance.CreateMetaInstanceRequest
	(*CreateMetaInstanceResponse)(nil),       // 1: proto.mrds.ledger.metainstance.CreateMetaInstanceResponse
	(*UpdateMetaInstanceStatusRequest)(nil),  // 2: proto.mrds.ledger.metainstance.UpdateMetaInstanceStatusRequest
	(*UpdateMetaInstanceResponse)(nil),       // 3: proto.mrds.ledger.metainstance.UpdateMetaInstanceResponse
	(*GetMetaInstanceByMetadataRequest)(nil), // 4: proto.mrds.ledger.metainstance.GetMetaInstanceByMetadataRequest
	(*GetMetaInstanceByNameRequest)(nil),     // 5: proto.mrds.ledger.metainstance.GetMetaInstanceByNameRequest
	(*GetMetaInstanceResponse)(nil),          // 6: proto.mrds.ledger.metainstance.GetMetaInstanceResponse
	(*ListMetaInstanceRequest)(nil),          // 7: proto.mrds.ledger.metainstance.ListMetaInstanceRequest
	(*ListMetaInstanceResponse)(nil),         // 8: proto.mrds.ledger.metainstance.ListMetaInstanceResponse
	(*DeleteMetaInstanceRequest)(nil),        // 9: proto.mrds.ledger.metainstance.DeleteMetaInstanceRequest
	(*DeleteMetaInstanceResponse)(nil),       // 10: proto.mrds.ledger.metainstance.DeleteMetaInstanceResponse
	(*MetaInstance)(nil),                     // 11: proto.mrds.ledger.metainstance.MetaInstance
	(*Metadata)(nil),                         // 12: proto.mrds.core.Metadata
	(*MetaInstanceStatus)(nil),               // 13: proto.mrds.ledger.metainstance.MetaInstanceStatus
	(MetaInstanceState)(0),                   // 14: proto.mrds.ledger.metainstance.MetaInstanceState
}
var file_metainstance_service_proto_depIdxs = []int32{
	11, // 0: proto.mrds.ledger.metainstance.CreateMetaInstanceResponse.record:type_name -> proto.mrds.ledger.metainstance.MetaInstance
	12, // 1: proto.mrds.ledger.metainstance.UpdateMetaInstanceStatusRequest.metadata:type_name -> proto.mrds.core.Metadata
	13, // 2: proto.mrds.ledger.metainstance.UpdateMetaInstanceStatusRequest.status:type_name -> proto.mrds.ledger.metainstance.MetaInstanceStatus
	11, // 3: proto.mrds.ledger.metainstance.UpdateMetaInstanceResponse.record:type_name -> proto.mrds.ledger.metainstance.MetaInstance
	12, // 4: proto.mrds.ledger.metainstance.GetMetaInstanceByMetadataRequest.metadata:type_name -> proto.mrds.core.Metadata
	11, // 5: proto.mrds.ledger.metainstance.GetMetaInstanceResponse.record:type_name -> proto.mrds.ledger.metainstance.MetaInstance
	14, // 6: proto.mrds.ledger.metainstance.ListMetaInstanceRequest.state_in:type_name -> proto.mrds.ledger.metainstance.MetaInstanceState
	14, // 7: proto.mrds.ledger.metainstance.ListMetaInstanceRequest.state_not_in:type_name -> proto.mrds.ledger.metainstance.MetaInstanceState
	11, // 8: proto.mrds.ledger.metainstance.ListMetaInstanceResponse.records:type_name -> proto.mrds.ledger.metainstance.MetaInstance
	12, // 9: proto.mrds.ledger.metainstance.DeleteMetaInstanceRequest.metadata:type_name -> proto.mrds.core.Metadata
	0,  // 10: proto.mrds.ledger.metainstance.MetaInstances.Create:input_type -> proto.mrds.ledger.metainstance.CreateMetaInstanceRequest
	4,  // 11: proto.mrds.ledger.metainstance.MetaInstances.GetByMetadata:input_type -> proto.mrds.ledger.metainstance.GetMetaInstanceByMetadataRequest
	5,  // 12: proto.mrds.ledger.metainstance.MetaInstances.GetByName:input_type -> proto.mrds.ledger.metainstance.GetMetaInstanceByNameRequest
	2,  // 13: proto.mrds.ledger.metainstance.MetaInstances.UpdateStatus:input_type -> proto.mrds.ledger.metainstance.UpdateMetaInstanceStatusRequest
	7,  // 14: proto.mrds.ledger.metainstance.MetaInstances.List:input_type -> proto.mrds.ledger.metainstance.ListMetaInstanceRequest
	9,  // 15: proto.mrds.ledger.metainstance.MetaInstances.Delete:input_type -> proto.mrds.ledger.metainstance.DeleteMetaInstanceRequest
	1,  // 16: proto.mrds.ledger.metainstance.MetaInstances.Create:output_type -> proto.mrds.ledger.metainstance.CreateMetaInstanceResponse
	6,  // 17: proto.mrds.ledger.metainstance.MetaInstances.GetByMetadata:output_type -> proto.mrds.ledger.metainstance.GetMetaInstanceResponse
	6,  // 18: proto.mrds.ledger.metainstance.MetaInstances.GetByName:output_type -> proto.mrds.ledger.metainstance.GetMetaInstanceResponse
	3,  // 19: proto.mrds.ledger.metainstance.MetaInstances.UpdateStatus:output_type -> proto.mrds.ledger.metainstance.UpdateMetaInstanceResponse
	8,  // 20: proto.mrds.ledger.metainstance.MetaInstances.List:output_type -> proto.mrds.ledger.metainstance.ListMetaInstanceResponse
	10, // 21: proto.mrds.ledger.metainstance.MetaInstances.Delete:output_type -> proto.mrds.ledger.metainstance.DeleteMetaInstanceResponse
	16, // [16:22] is the sub-list for method output_type
	10, // [10:16] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_metainstance_service_proto_init() }
func file_metainstance_service_proto_init() {
	if File_metainstance_service_proto != nil {
		return
	}
	file_metadata_proto_init()
	file_metainstance_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metainstance_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metainstance_service_proto_goTypes,
		DependencyIndexes: file_metainstance_service_proto_depIdxs,
		MessageInfos:      file_metainstance_service_proto_msgTypes,
	}.Build()
	File_metainstance_service_proto = out.File
	file_metainstance_service_proto_rawDesc = nil
	file_metainstance_service_proto_goTypes = nil
	file_metainstance_service_proto_depIdxs = nil
}