// 这个就是protobuf的中间文件

// 指定的当前proto语法的版本，有2和3

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.11
// source: server.proto

// 指定等会文件生成出来的package

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// 定义request model
type PlumelogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic   string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Env     string `protobuf:"bytes,2,opt,name=env,proto3" json:"env,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"` // 1代表顺序
}

func (x *PlumelogRequest) Reset() {
	*x = PlumelogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumelogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumelogRequest) ProtoMessage() {}

func (x *PlumelogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumelogRequest.ProtoReflect.Descriptor instead.
func (*PlumelogRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{0}
}

func (x *PlumelogRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *PlumelogRequest) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *PlumelogRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// 定义response model
type PlumelogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // 1代表顺序
}

func (x *PlumelogResponse) Reset() {
	*x = PlumelogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumelogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumelogResponse) ProtoMessage() {}

func (x *PlumelogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumelogResponse.ProtoReflect.Descriptor instead.
func (*PlumelogResponse) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{1}
}

func (x *PlumelogResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PlumelogInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DateTime   string `protobuf:"bytes,1,opt,name=dateTime,proto3" json:"dateTime,omitempty"`
	TraceId    string `protobuf:"bytes,2,opt,name=traceId,proto3" json:"traceId,omitempty"`
	Method     string `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	AppName    string `protobuf:"bytes,4,opt,name=appName,proto3" json:"appName,omitempty"`
	ServerName string `protobuf:"bytes,5,opt,name=serverName,proto3" json:"serverName,omitempty"`
	ClassName  string `protobuf:"bytes,6,opt,name=className,proto3" json:"className,omitempty"`
	Env        string `protobuf:"bytes,7,opt,name=env,proto3" json:"env,omitempty"`
	Content    string `protobuf:"bytes,8,opt,name=content,proto3" json:"content,omitempty"`
	ThreadName string `protobuf:"bytes,9,opt,name=threadName,proto3" json:"threadName,omitempty"`
	Url        string `protobuf:"bytes,10,opt,name=url,proto3" json:"url,omitempty"`
	SpanId     string `protobuf:"bytes,11,opt,name=spanId,proto3" json:"spanId,omitempty"`
	DtTime     int64  `protobuf:"varint,12,opt,name=dtTime,proto3" json:"dtTime,omitempty"`
	CostTime   int64  `protobuf:"varint,13,opt,name=costTime,proto3" json:"costTime,omitempty"`
	LogLevel   string `protobuf:"bytes,14,opt,name=logLevel,proto3" json:"logLevel,omitempty"`
	Seq        int32  `protobuf:"varint,15,opt,name=seq,proto3" json:"seq,omitempty"`
}

func (x *PlumelogInfo) Reset() {
	*x = PlumelogInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumelogInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumelogInfo) ProtoMessage() {}

func (x *PlumelogInfo) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumelogInfo.ProtoReflect.Descriptor instead.
func (*PlumelogInfo) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{2}
}

func (x *PlumelogInfo) GetDateTime() string {
	if x != nil {
		return x.DateTime
	}
	return ""
}

func (x *PlumelogInfo) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *PlumelogInfo) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *PlumelogInfo) GetAppName() string {
	if x != nil {
		return x.AppName
	}
	return ""
}

func (x *PlumelogInfo) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *PlumelogInfo) GetClassName() string {
	if x != nil {
		return x.ClassName
	}
	return ""
}

func (x *PlumelogInfo) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *PlumelogInfo) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *PlumelogInfo) GetThreadName() string {
	if x != nil {
		return x.ThreadName
	}
	return ""
}

func (x *PlumelogInfo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *PlumelogInfo) GetSpanId() string {
	if x != nil {
		return x.SpanId
	}
	return ""
}

func (x *PlumelogInfo) GetDtTime() int64 {
	if x != nil {
		return x.DtTime
	}
	return 0
}

func (x *PlumelogInfo) GetCostTime() int64 {
	if x != nil {
		return x.CostTime
	}
	return 0
}

func (x *PlumelogInfo) GetLogLevel() string {
	if x != nil {
		return x.LogLevel
	}
	return ""
}

func (x *PlumelogInfo) GetSeq() int32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

type PlumelogServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=serviceName,proto3" json:"serviceName,omitempty"` // 1代表顺序
}

func (x *PlumelogServiceRequest) Reset() {
	*x = PlumelogServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumelogServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumelogServiceRequest) ProtoMessage() {}

func (x *PlumelogServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumelogServiceRequest.ProtoReflect.Descriptor instead.
func (*PlumelogServiceRequest) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{3}
}

func (x *PlumelogServiceRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

var File_server_proto protoreflect.FileDescriptor

var file_server_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x53, 0x0a, 0x0f, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c,
	0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70,
	0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12,
	0x10, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e,
	0x76, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2c, 0x0a, 0x10, 0x50,
	0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x8c, 0x03, 0x0a, 0x0c, 0x50, 0x6c,
	0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x70, 0x70, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65,
	0x6e, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x64, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x6f, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x63, 0x6f, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f,
	0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x73, 0x65, 0x71, 0x22, 0x3a, 0x0a, 0x16, 0x50, 0x6c, 0x75, 0x6d,
	0x65, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x32, 0xd9, 0x01, 0x0a, 0x0f, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f,
	0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64,
	0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x12, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x65,
	0x6c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0a, 0x50,
	0x75, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1e, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0d, 0x57, 0x72, 0x69, 0x74, 0x65, 0x50, 0x6c, 0x75, 0x6d,
	0x65, 0x4c, 0x6f, 0x67, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x6c,
	0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x49, 0x6e, 0x66, 0x6f,
	0x42, 0x15, 0x5a, 0x13, 0x70, 0x6c, 0x75, 0x6d, 0x65, 0x6c, 0x6f, 0x67, 0x2f, 0x72, 0x70, 0x63,
	0x3b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_rawDescOnce sync.Once
	file_server_proto_rawDescData = file_server_proto_rawDesc
)

func file_server_proto_rawDescGZIP() []byte {
	file_server_proto_rawDescOnce.Do(func() {
		file_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_rawDescData)
	})
	return file_server_proto_rawDescData
}

var file_server_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_server_proto_goTypes = []interface{}{
	(*PlumelogRequest)(nil),        // 0: server.PlumelogRequest
	(*PlumelogResponse)(nil),       // 1: server.PlumelogResponse
	(*PlumelogInfo)(nil),           // 2: server.PlumelogInfo
	(*PlumelogServiceRequest)(nil), // 3: server.PlumelogServiceRequest
}
var file_server_proto_depIdxs = []int32{
	0, // 0: server.PlumelogService.SendPlumelog:input_type -> server.PlumelogRequest
	3, // 1: server.PlumelogService.PutService:input_type -> server.PlumelogServiceRequest
	2, // 2: server.PlumelogService.WritePlumeLog:input_type -> server.PlumelogInfo
	1, // 3: server.PlumelogService.SendPlumelog:output_type -> server.PlumelogResponse
	1, // 4: server.PlumelogService.PutService:output_type -> server.PlumelogResponse
	2, // 5: server.PlumelogService.WritePlumeLog:output_type -> server.PlumelogInfo
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_server_proto_init() }
func file_server_proto_init() {
	if File_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumelogRequest); i {
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
		file_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumelogResponse); i {
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
		file_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumelogInfo); i {
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
		file_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumelogServiceRequest); i {
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
			RawDescriptor: file_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_goTypes,
		DependencyIndexes: file_server_proto_depIdxs,
		MessageInfos:      file_server_proto_msgTypes,
	}.Build()
	File_server_proto = out.File
	file_server_proto_rawDesc = nil
	file_server_proto_goTypes = nil
	file_server_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PlumelogServiceClient is the client API for PlumelogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PlumelogServiceClient interface {
	// 定义方法
	SendPlumelog(ctx context.Context, in *PlumelogRequest, opts ...grpc.CallOption) (*PlumelogResponse, error)
	PutService(ctx context.Context, in *PlumelogServiceRequest, opts ...grpc.CallOption) (*PlumelogResponse, error)
	WritePlumeLog(ctx context.Context, in *PlumelogInfo, opts ...grpc.CallOption) (*PlumelogInfo, error)
}

type plumelogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlumelogServiceClient(cc grpc.ClientConnInterface) PlumelogServiceClient {
	return &plumelogServiceClient{cc}
}

func (c *plumelogServiceClient) SendPlumelog(ctx context.Context, in *PlumelogRequest, opts ...grpc.CallOption) (*PlumelogResponse, error) {
	out := new(PlumelogResponse)
	err := c.cc.Invoke(ctx, "/server.PlumelogService/SendPlumelog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plumelogServiceClient) PutService(ctx context.Context, in *PlumelogServiceRequest, opts ...grpc.CallOption) (*PlumelogResponse, error) {
	out := new(PlumelogResponse)
	err := c.cc.Invoke(ctx, "/server.PlumelogService/PutService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plumelogServiceClient) WritePlumeLog(ctx context.Context, in *PlumelogInfo, opts ...grpc.CallOption) (*PlumelogInfo, error) {
	out := new(PlumelogInfo)
	err := c.cc.Invoke(ctx, "/server.PlumelogService/WritePlumeLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlumelogServiceServer is the server API for PlumelogService service.
type PlumelogServiceServer interface {
	// 定义方法
	SendPlumelog(context.Context, *PlumelogRequest) (*PlumelogResponse, error)
	PutService(context.Context, *PlumelogServiceRequest) (*PlumelogResponse, error)
	WritePlumeLog(context.Context, *PlumelogInfo) (*PlumelogInfo, error)
}

// UnimplementedPlumelogServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPlumelogServiceServer struct {
}

func (*UnimplementedPlumelogServiceServer) SendPlumelog(context.Context, *PlumelogRequest) (*PlumelogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPlumelog not implemented")
}
func (*UnimplementedPlumelogServiceServer) PutService(context.Context, *PlumelogServiceRequest) (*PlumelogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutService not implemented")
}
func (*UnimplementedPlumelogServiceServer) WritePlumeLog(context.Context, *PlumelogInfo) (*PlumelogInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WritePlumeLog not implemented")
}

func RegisterPlumelogServiceServer(s *grpc.Server, srv PlumelogServiceServer) {
	s.RegisterService(&_PlumelogService_serviceDesc, srv)
}

func _PlumelogService_SendPlumelog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumelogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumelogServiceServer).SendPlumelog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.PlumelogService/SendPlumelog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumelogServiceServer).SendPlumelog(ctx, req.(*PlumelogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlumelogService_PutService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumelogServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumelogServiceServer).PutService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.PlumelogService/PutService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumelogServiceServer).PutService(ctx, req.(*PlumelogServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlumelogService_WritePlumeLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumelogInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumelogServiceServer).WritePlumeLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.PlumelogService/WritePlumeLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumelogServiceServer).WritePlumeLog(ctx, req.(*PlumelogInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _PlumelogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "server.PlumelogService",
	HandlerType: (*PlumelogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendPlumelog",
			Handler:    _PlumelogService_SendPlumelog_Handler,
		},
		{
			MethodName: "PutService",
			Handler:    _PlumelogService_PutService_Handler,
		},
		{
			MethodName: "WritePlumeLog",
			Handler:    _PlumelogService_WritePlumeLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
