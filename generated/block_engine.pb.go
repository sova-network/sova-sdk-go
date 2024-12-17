// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.3
// source: block_engine.proto

package types

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

type StreamMempoolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StreamMempoolResponse) Reset() {
	*x = StreamMempoolResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_engine_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamMempoolResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamMempoolResponse) ProtoMessage() {}

func (x *StreamMempoolResponse) ProtoReflect() protoreflect.Message {
	mi := &file_block_engine_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamMempoolResponse.ProtoReflect.Descriptor instead.
func (*StreamMempoolResponse) Descriptor() ([]byte, []int) {
	return file_block_engine_proto_rawDescGZIP(), []int{0}
}

type SubscribeBundlesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SubscribeBundlesRequest) Reset() {
	*x = SubscribeBundlesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_engine_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeBundlesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeBundlesRequest) ProtoMessage() {}

func (x *SubscribeBundlesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_block_engine_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeBundlesRequest.ProtoReflect.Descriptor instead.
func (*SubscribeBundlesRequest) Descriptor() ([]byte, []int) {
	return file_block_engine_proto_rawDescGZIP(), []int{1}
}

var File_block_engine_proto protoreflect.FileDescriptor

var file_block_engine_proto_rawDesc = []byte{
	0x0a, 0x12, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x1a, 0x09, 0x64, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x17, 0x0a,
	0x15, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x6d, 0x70, 0x6f, 0x6f, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x19, 0x0a, 0x17, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x32, 0xb0, 0x01, 0x0a, 0x14, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x4c, 0x0a, 0x0d, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x6d, 0x70, 0x6f, 0x6f, 0x6c, 0x12, 0x12, 0x2e, 0x64, 0x74,
	0x6f, 0x2e, 0x4d, 0x65, 0x6d, 0x70, 0x6f, 0x6f, 0x6c, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x1a,
	0x23, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x6d, 0x70, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x4a, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x12, 0x25, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x64, 0x74, 0x6f, 0x2e, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x64, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_engine_proto_rawDescOnce sync.Once
	file_block_engine_proto_rawDescData = file_block_engine_proto_rawDesc
)

func file_block_engine_proto_rawDescGZIP() []byte {
	file_block_engine_proto_rawDescOnce.Do(func() {
		file_block_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_engine_proto_rawDescData)
	})
	return file_block_engine_proto_rawDescData
}

var file_block_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_block_engine_proto_goTypes = []any{
	(*StreamMempoolResponse)(nil),   // 0: block_engine.StreamMempoolResponse
	(*SubscribeBundlesRequest)(nil), // 1: block_engine.SubscribeBundlesRequest
	(*MempoolPacket)(nil),           // 2: dto.MempoolPacket
	(*Bundle)(nil),                  // 3: dto.Bundle
}
var file_block_engine_proto_depIdxs = []int32{
	2, // 0: block_engine.BlockEngineValidator.StreamMempool:input_type -> dto.MempoolPacket
	1, // 1: block_engine.BlockEngineValidator.SubscribeBundles:input_type -> block_engine.SubscribeBundlesRequest
	0, // 2: block_engine.BlockEngineValidator.StreamMempool:output_type -> block_engine.StreamMempoolResponse
	3, // 3: block_engine.BlockEngineValidator.SubscribeBundles:output_type -> dto.Bundle
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_block_engine_proto_init() }
func file_block_engine_proto_init() {
	if File_block_engine_proto != nil {
		return
	}
	file_dto_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_block_engine_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*StreamMempoolResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_engine_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*SubscribeBundlesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_block_engine_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_block_engine_proto_goTypes,
		DependencyIndexes: file_block_engine_proto_depIdxs,
		MessageInfos:      file_block_engine_proto_msgTypes,
	}.Build()
	File_block_engine_proto = out.File
	file_block_engine_proto_rawDesc = nil
	file_block_engine_proto_goTypes = nil
	file_block_engine_proto_depIdxs = nil
}
