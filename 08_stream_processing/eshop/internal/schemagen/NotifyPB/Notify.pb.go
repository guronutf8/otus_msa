// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: schema/Notify.proto

package NotifyPB

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

type SendNotifyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User   string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Result bool   `protobuf:"varint,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *SendNotifyRequest) Reset() {
	*x = SendNotifyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_Notify_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendNotifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendNotifyRequest) ProtoMessage() {}

func (x *SendNotifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_schema_Notify_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendNotifyRequest.ProtoReflect.Descriptor instead.
func (*SendNotifyRequest) Descriptor() ([]byte, []int) {
	return file_schema_Notify_proto_rawDescGZIP(), []int{0}
}

func (x *SendNotifyRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *SendNotifyRequest) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

type SendNotifyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendNotifyResponse) Reset() {
	*x = SendNotifyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_Notify_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendNotifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendNotifyResponse) ProtoMessage() {}

func (x *SendNotifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_schema_Notify_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendNotifyResponse.ProtoReflect.Descriptor instead.
func (*SendNotifyResponse) Descriptor() ([]byte, []int) {
	return file_schema_Notify_proto_rawDescGZIP(), []int{1}
}

var File_schema_Notify_proto protoreflect.FileDescriptor

var file_schema_Notify_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x22, 0x3f, 0x0a,
	0x11, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x14,
	0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x4f, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x45,
	0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x19, 0x2e, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x50, 0x42, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_schema_Notify_proto_rawDescOnce sync.Once
	file_schema_Notify_proto_rawDescData = file_schema_Notify_proto_rawDesc
)

func file_schema_Notify_proto_rawDescGZIP() []byte {
	file_schema_Notify_proto_rawDescOnce.Do(func() {
		file_schema_Notify_proto_rawDescData = protoimpl.X.CompressGZIP(file_schema_Notify_proto_rawDescData)
	})
	return file_schema_Notify_proto_rawDescData
}

var file_schema_Notify_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_schema_Notify_proto_goTypes = []interface{}{
	(*SendNotifyRequest)(nil),  // 0: Notify.SendNotifyRequest
	(*SendNotifyResponse)(nil), // 1: Notify.SendNotifyResponse
}
var file_schema_Notify_proto_depIdxs = []int32{
	0, // 0: Notify.Notify.SendNotify:input_type -> Notify.SendNotifyRequest
	1, // 1: Notify.Notify.SendNotify:output_type -> Notify.SendNotifyResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_schema_Notify_proto_init() }
func file_schema_Notify_proto_init() {
	if File_schema_Notify_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_schema_Notify_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendNotifyRequest); i {
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
		file_schema_Notify_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendNotifyResponse); i {
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
			RawDescriptor: file_schema_Notify_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_schema_Notify_proto_goTypes,
		DependencyIndexes: file_schema_Notify_proto_depIdxs,
		MessageInfos:      file_schema_Notify_proto_msgTypes,
	}.Build()
	File_schema_Notify_proto = out.File
	file_schema_Notify_proto_rawDesc = nil
	file_schema_Notify_proto_goTypes = nil
	file_schema_Notify_proto_depIdxs = nil
}
