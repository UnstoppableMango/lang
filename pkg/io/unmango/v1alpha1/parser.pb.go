// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: io/unmango/v1alpha1/parser.proto

package unmangov1alpha1

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

type ParseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ParseRequest) Reset() {
	*x = ParseRequest{}
	mi := &file_io_unmango_v1alpha1_parser_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseRequest) ProtoMessage() {}

func (x *ParseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_io_unmango_v1alpha1_parser_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseRequest.ProtoReflect.Descriptor instead.
func (*ParseRequest) Descriptor() ([]byte, []int) {
	return file_io_unmango_v1alpha1_parser_proto_rawDescGZIP(), []int{0}
}

type ParseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ParseResponse) Reset() {
	*x = ParseResponse{}
	mi := &file_io_unmango_v1alpha1_parser_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseResponse) ProtoMessage() {}

func (x *ParseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_io_unmango_v1alpha1_parser_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseResponse.ProtoReflect.Descriptor instead.
func (*ParseResponse) Descriptor() ([]byte, []int) {
	return file_io_unmango_v1alpha1_parser_proto_rawDescGZIP(), []int{1}
}

var File_io_unmango_v1alpha1_parser_proto protoreflect.FileDescriptor

var file_io_unmango_v1alpha1_parser_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6f, 0x2f, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x13, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x61, 0x72, 0x73, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0f, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x73, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x5f, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x73,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x05, 0x50, 0x61, 0x72,
	0x73, 0x65, 0x12, 0x21, 0x2e, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e,
	0x67, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x50, 0x61, 0x72, 0x73,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xde, 0x01, 0x0a, 0x17, 0x63, 0x6f,
	0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0b, 0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x75, 0x6e, 0x73, 0x74, 0x6f, 0x70, 0x70, 0x61, 0x62, 0x6c, 0x65, 0x6d, 0x61, 0x6e, 0x67,
	0x6f, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x69, 0x6f, 0x2f, 0x75, 0x6e,
	0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x75,
	0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02,
	0x03, 0x49, 0x55, 0x58, 0xaa, 0x02, 0x13, 0x49, 0x6f, 0x2e, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67,
	0x6f, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x13, 0x49, 0x6f, 0x5c,
	0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0xe2, 0x02, 0x1f, 0x49, 0x6f, 0x5c, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x15, 0x49, 0x6f, 0x3a, 0x3a, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f,
	0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_io_unmango_v1alpha1_parser_proto_rawDescOnce sync.Once
	file_io_unmango_v1alpha1_parser_proto_rawDescData = file_io_unmango_v1alpha1_parser_proto_rawDesc
)

func file_io_unmango_v1alpha1_parser_proto_rawDescGZIP() []byte {
	file_io_unmango_v1alpha1_parser_proto_rawDescOnce.Do(func() {
		file_io_unmango_v1alpha1_parser_proto_rawDescData = protoimpl.X.CompressGZIP(file_io_unmango_v1alpha1_parser_proto_rawDescData)
	})
	return file_io_unmango_v1alpha1_parser_proto_rawDescData
}

var file_io_unmango_v1alpha1_parser_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_io_unmango_v1alpha1_parser_proto_goTypes = []any{
	(*ParseRequest)(nil),  // 0: io.unmango.v1alpha1.ParseRequest
	(*ParseResponse)(nil), // 1: io.unmango.v1alpha1.ParseResponse
}
var file_io_unmango_v1alpha1_parser_proto_depIdxs = []int32{
	0, // 0: io.unmango.v1alpha1.ParserService.Parse:input_type -> io.unmango.v1alpha1.ParseRequest
	1, // 1: io.unmango.v1alpha1.ParserService.Parse:output_type -> io.unmango.v1alpha1.ParseResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_io_unmango_v1alpha1_parser_proto_init() }
func file_io_unmango_v1alpha1_parser_proto_init() {
	if File_io_unmango_v1alpha1_parser_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_io_unmango_v1alpha1_parser_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_io_unmango_v1alpha1_parser_proto_goTypes,
		DependencyIndexes: file_io_unmango_v1alpha1_parser_proto_depIdxs,
		MessageInfos:      file_io_unmango_v1alpha1_parser_proto_msgTypes,
	}.Build()
	File_io_unmango_v1alpha1_parser_proto = out.File
	file_io_unmango_v1alpha1_parser_proto_rawDesc = nil
	file_io_unmango_v1alpha1_parser_proto_goTypes = nil
	file_io_unmango_v1alpha1_parser_proto_depIdxs = nil
}
