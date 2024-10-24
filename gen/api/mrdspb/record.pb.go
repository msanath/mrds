// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: record.proto

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

// Enum to represent the RecordState
type RecordState int32

const (
	RecordState_RecordState_UNKNOWN  RecordState = 0
	RecordState_RecordState_ACTIVE   RecordState = 1
	RecordState_RecordState_INACTIVE RecordState = 2
)

// Enum value maps for RecordState.
var (
	RecordState_name = map[int32]string{
		0: "RecordState_UNKNOWN",
		1: "RecordState_ACTIVE",
		2: "RecordState_INACTIVE",
	}
	RecordState_value = map[string]int32{
		"RecordState_UNKNOWN":  0,
		"RecordState_ACTIVE":   1,
		"RecordState_INACTIVE": 2,
	}
)

func (x RecordState) Enum() *RecordState {
	p := new(RecordState)
	*p = x
	return p
}

func (x RecordState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecordState) Descriptor() protoreflect.EnumDescriptor {
	return file_record_proto_enumTypes[0].Descriptor()
}

func (RecordState) Type() protoreflect.EnumType {
	return &file_record_proto_enumTypes[0]
}

func (x RecordState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecordState.Descriptor instead.
func (RecordState) EnumDescriptor() ([]byte, []int) {
	return file_record_proto_rawDescGZIP(), []int{0}
}

// Message representing the Record
type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Metadata is the metadata that identifies the Record.
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Name is the name of the Record.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Status represents the current status of the Record.
	Status *RecordStatus `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	mi := &file_record_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_record_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_record_proto_rawDescGZIP(), []int{0}
}

func (x *Record) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Record) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Record) GetStatus() *RecordStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// Message representing the Status of a resource.
type RecordStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// State is the discrete condition of the resource.
	State RecordState `protobuf:"varint,1,opt,name=state,proto3,enum=proto.mrds.ledger.record.RecordState" json:"state,omitempty"`
	// Message is a human-readable description of the resource's state.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RecordStatus) Reset() {
	*x = RecordStatus{}
	mi := &file_record_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecordStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordStatus) ProtoMessage() {}

func (x *RecordStatus) ProtoReflect() protoreflect.Message {
	mi := &file_record_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordStatus.ProtoReflect.Descriptor instead.
func (*RecordStatus) Descriptor() ([]byte, []int) {
	return file_record_proto_rawDescGZIP(), []int{1}
}

func (x *RecordStatus) GetState() RecordState {
	if x != nil {
		return x.State
	}
	return RecordState_RecordState_UNKNOWN
}

func (x *RecordStatus) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_record_proto protoreflect.FileDescriptor

var file_record_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x1a, 0x0e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72,
	0x64, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3e,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x65,
	0x0a, 0x0c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3b,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x72, 0x64, 0x73, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x58, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x16, 0x0a,
	0x12, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x41, 0x43, 0x54,
	0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x5f, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x42,
	0x0d, 0x5a, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x72, 0x64, 0x73, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_record_proto_rawDescOnce sync.Once
	file_record_proto_rawDescData = file_record_proto_rawDesc
)

func file_record_proto_rawDescGZIP() []byte {
	file_record_proto_rawDescOnce.Do(func() {
		file_record_proto_rawDescData = protoimpl.X.CompressGZIP(file_record_proto_rawDescData)
	})
	return file_record_proto_rawDescData
}

var file_record_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_record_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_record_proto_goTypes = []any{
	(RecordState)(0),     // 0: proto.mrds.ledger.record.RecordState
	(*Record)(nil),       // 1: proto.mrds.ledger.record.Record
	(*RecordStatus)(nil), // 2: proto.mrds.ledger.record.RecordStatus
	(*Metadata)(nil),     // 3: proto.mrds.core.Metadata
}
var file_record_proto_depIdxs = []int32{
	3, // 0: proto.mrds.ledger.record.Record.metadata:type_name -> proto.mrds.core.Metadata
	2, // 1: proto.mrds.ledger.record.Record.status:type_name -> proto.mrds.ledger.record.RecordStatus
	0, // 2: proto.mrds.ledger.record.RecordStatus.state:type_name -> proto.mrds.ledger.record.RecordState
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_record_proto_init() }
func file_record_proto_init() {
	if File_record_proto != nil {
		return
	}
	file_metadata_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_record_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_record_proto_goTypes,
		DependencyIndexes: file_record_proto_depIdxs,
		EnumInfos:         file_record_proto_enumTypes,
		MessageInfos:      file_record_proto_msgTypes,
	}.Build()
	File_record_proto = out.File
	file_record_proto_rawDesc = nil
	file_record_proto_goTypes = nil
	file_record_proto_depIdxs = nil
}
