// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: metainstance.proto

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

// Enum to represent the MetaInstanceState
type MetaInstanceState int32

const (
	MetaInstanceState_MetaInstanceState_UNKNOWN  MetaInstanceState = 0
	MetaInstanceState_MetaInstanceState_PENDING  MetaInstanceState = 1
	MetaInstanceState_MetaInstanceState_ACTIVE   MetaInstanceState = 2
	MetaInstanceState_MetaInstanceState_INACTIVE MetaInstanceState = 3
)

// Enum value maps for MetaInstanceState.
var (
	MetaInstanceState_name = map[int32]string{
		0: "MetaInstanceState_UNKNOWN",
		1: "MetaInstanceState_PENDING",
		2: "MetaInstanceState_ACTIVE",
		3: "MetaInstanceState_INACTIVE",
	}
	MetaInstanceState_value = map[string]int32{
		"MetaInstanceState_UNKNOWN":  0,
		"MetaInstanceState_PENDING":  1,
		"MetaInstanceState_ACTIVE":   2,
		"MetaInstanceState_INACTIVE": 3,
	}
)

func (x MetaInstanceState) Enum() *MetaInstanceState {
	p := new(MetaInstanceState)
	*p = x
	return p
}

func (x MetaInstanceState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetaInstanceState) Descriptor() protoreflect.EnumDescriptor {
	return file_metainstance_proto_enumTypes[0].Descriptor()
}

func (MetaInstanceState) Type() protoreflect.EnumType {
	return &file_metainstance_proto_enumTypes[0]
}

func (x MetaInstanceState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetaInstanceState.Descriptor instead.
func (MetaInstanceState) EnumDescriptor() ([]byte, []int) {
	return file_metainstance_proto_rawDescGZIP(), []int{0}
}

// Message representing the MetaInstance
type MetaInstance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Metadata is the metadata that identifies the MetaInstance.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Name is the name of the MetaInstance.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Status represents the current status of the MetaInstance.
	Status *MetaInstanceStatus `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *MetaInstance) Reset() {
	*x = MetaInstance{}
	mi := &file_metainstance_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetaInstance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaInstance) ProtoMessage() {}

func (x *MetaInstance) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaInstance.ProtoReflect.Descriptor instead.
func (*MetaInstance) Descriptor() ([]byte, []int) {
	return file_metainstance_proto_rawDescGZIP(), []int{0}
}

func (x *MetaInstance) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *MetaInstance) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MetaInstance) GetStatus() *MetaInstanceStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// Message representing the Status of a resource.
type MetaInstanceStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// State is the discrete condition of the resource.
	State MetaInstanceState `protobuf:"varint,1,opt,name=state,proto3,enum=proto.mrds.ledger.metainstance.MetaInstanceState" json:"state,omitempty"`
	// Message is a human-readable description of the resource's state.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MetaInstanceStatus) Reset() {
	*x = MetaInstanceStatus{}
	mi := &file_metainstance_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetaInstanceStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaInstanceStatus) ProtoMessage() {}

func (x *MetaInstanceStatus) ProtoReflect() protoreflect.Message {
	mi := &file_metainstance_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaInstanceStatus.ProtoReflect.Descriptor instead.
func (*MetaInstanceStatus) Descriptor() ([]byte, []int) {
	return file_metainstance_proto_rawDescGZIP(), []int{1}
}

func (x *MetaInstanceStatus) GetState() MetaInstanceState {
	if x != nil {
		return x.State
	}
	return MetaInstanceState_MetaInstanceState_UNKNOWN
}

func (x *MetaInstanceStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_metainstance_proto protoreflect.FileDescriptor

var file_metainstance_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73,
	0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x1a, 0x0e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa5, 0x01, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6d, 0x72, 0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x4a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x77, 0x0a, 0x12,
	0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x47, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x8f, 0x01, 0x0a, 0x11, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x4d, 0x65,
	0x74, 0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x4d, 0x65, 0x74,
	0x61, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x4d, 0x65, 0x74, 0x61, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x49, 0x4e, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03, 0x42, 0x0d, 0x5a, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x6d, 0x72, 0x64, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metainstance_proto_rawDescOnce sync.Once
	file_metainstance_proto_rawDescData = file_metainstance_proto_rawDesc
)

func file_metainstance_proto_rawDescGZIP() []byte {
	file_metainstance_proto_rawDescOnce.Do(func() {
		file_metainstance_proto_rawDescData = protoimpl.X.CompressGZIP(file_metainstance_proto_rawDescData)
	})
	return file_metainstance_proto_rawDescData
}

var file_metainstance_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_metainstance_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_metainstance_proto_goTypes = []any{
	(MetaInstanceState)(0),     // 0: proto.mrds.ledger.metainstance.MetaInstanceState
	(*MetaInstance)(nil),       // 1: proto.mrds.ledger.metainstance.MetaInstance
	(*MetaInstanceStatus)(nil), // 2: proto.mrds.ledger.metainstance.MetaInstanceStatus
	(*Metadata)(nil),           // 3: proto.mrds.core.Metadata
}
var file_metainstance_proto_depIdxs = []int32{
	3, // 0: proto.mrds.ledger.metainstance.MetaInstance.metadata:type_name -> proto.mrds.core.Metadata
	2, // 1: proto.mrds.ledger.metainstance.MetaInstance.status:type_name -> proto.mrds.ledger.metainstance.MetaInstanceStatus
	0, // 2: proto.mrds.ledger.metainstance.MetaInstanceStatus.state:type_name -> proto.mrds.ledger.metainstance.MetaInstanceState
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_metainstance_proto_init() }
func file_metainstance_proto_init() {
	if File_metainstance_proto != nil {
		return
	}
	file_metadata_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metainstance_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_metainstance_proto_goTypes,
		DependencyIndexes: file_metainstance_proto_depIdxs,
		EnumInfos:         file_metainstance_proto_enumTypes,
		MessageInfos:      file_metainstance_proto_msgTypes,
	}.Build()
	File_metainstance_proto = out.File
	file_metainstance_proto_rawDesc = nil
	file_metainstance_proto_goTypes = nil
	file_metainstance_proto_depIdxs = nil
}
