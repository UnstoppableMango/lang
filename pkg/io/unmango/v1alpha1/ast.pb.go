// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: io/unmango/v1alpha1/ast.proto

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

type String struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *String) Reset() {
	*x = String{}
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *String) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*String) ProtoMessage() {}

func (x *String) ProtoReflect() protoreflect.Message {
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use String.ProtoReflect.Descriptor instead.
func (*String) Descriptor() ([]byte, []int) {
	return file_io_unmango_v1alpha1_ast_proto_rawDescGZIP(), []int{0}
}

func (x *String) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*Node_String_
	Value isNode_Value `protobuf_oneof:"value"`
}

func (x *Node) Reset() {
	*x = Node{}
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_io_unmango_v1alpha1_ast_proto_rawDescGZIP(), []int{1}
}

func (m *Node) GetValue() isNode_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Node) GetString_() *String {
	if x, ok := x.GetValue().(*Node_String_); ok {
		return x.String_
	}
	return nil
}

type isNode_Value interface {
	isNode_Value()
}

type Node_String_ struct {
	String_ *String `protobuf:"bytes,1,opt,name=string,proto3,oneof"`
}

func (*Node_String_) isNode_Value() {}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_io_unmango_v1alpha1_ast_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_io_unmango_v1alpha1_ast_proto_rawDescGZIP(), []int{2}
}

func (x *File) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

var File_io_unmango_v1alpha1_ast_proto protoreflect.FileDescriptor

var file_io_unmango_v1alpha1_ast_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x6f, 0x2f, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x61, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x22, 0x1e, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x46, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x35, 0x0a, 0x06,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69,
	0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x37, 0x0a, 0x04,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x69, 0x6f, 0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05,
	0x6e, 0x6f, 0x64, 0x65, 0x73, 0x42, 0xdb, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6f,
	0x2e, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x42, 0x08, 0x41, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x48, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6e, 0x73, 0x74, 0x6f, 0x70,
	0x70, 0x61, 0x62, 0x6c, 0x65, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x6c, 0x61, 0x6e, 0x67, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x69, 0x6f, 0x2f, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x75, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x55, 0x58, 0xaa, 0x02, 0x13,
	0x49, 0x6f, 0x2e, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0xca, 0x02, 0x13, 0x49, 0x6f, 0x5c, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f,
	0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x1f, 0x49, 0x6f, 0x5c, 0x55,
	0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x15, 0x49, 0x6f,
	0x3a, 0x3a, 0x55, 0x6e, 0x6d, 0x61, 0x6e, 0x67, 0x6f, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_io_unmango_v1alpha1_ast_proto_rawDescOnce sync.Once
	file_io_unmango_v1alpha1_ast_proto_rawDescData = file_io_unmango_v1alpha1_ast_proto_rawDesc
)

func file_io_unmango_v1alpha1_ast_proto_rawDescGZIP() []byte {
	file_io_unmango_v1alpha1_ast_proto_rawDescOnce.Do(func() {
		file_io_unmango_v1alpha1_ast_proto_rawDescData = protoimpl.X.CompressGZIP(file_io_unmango_v1alpha1_ast_proto_rawDescData)
	})
	return file_io_unmango_v1alpha1_ast_proto_rawDescData
}

var file_io_unmango_v1alpha1_ast_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_io_unmango_v1alpha1_ast_proto_goTypes = []any{
	(*String)(nil), // 0: io.unmango.v1alpha1.String
	(*Node)(nil),   // 1: io.unmango.v1alpha1.Node
	(*File)(nil),   // 2: io.unmango.v1alpha1.File
}
var file_io_unmango_v1alpha1_ast_proto_depIdxs = []int32{
	0, // 0: io.unmango.v1alpha1.Node.string:type_name -> io.unmango.v1alpha1.String
	1, // 1: io.unmango.v1alpha1.File.nodes:type_name -> io.unmango.v1alpha1.Node
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_io_unmango_v1alpha1_ast_proto_init() }
func file_io_unmango_v1alpha1_ast_proto_init() {
	if File_io_unmango_v1alpha1_ast_proto != nil {
		return
	}
	file_io_unmango_v1alpha1_ast_proto_msgTypes[1].OneofWrappers = []any{
		(*Node_String_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_io_unmango_v1alpha1_ast_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_io_unmango_v1alpha1_ast_proto_goTypes,
		DependencyIndexes: file_io_unmango_v1alpha1_ast_proto_depIdxs,
		MessageInfos:      file_io_unmango_v1alpha1_ast_proto_msgTypes,
	}.Build()
	File_io_unmango_v1alpha1_ast_proto = out.File
	file_io_unmango_v1alpha1_ast_proto_rawDesc = nil
	file_io_unmango_v1alpha1_ast_proto_goTypes = nil
	file_io_unmango_v1alpha1_ast_proto_depIdxs = nil
}
