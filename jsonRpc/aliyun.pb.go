// Code generated by protoc-gen-go. DO NOT EDIT.
// source: aliyun.proto

package app

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SendRequest struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	ParamStr             string   `protobuf:"bytes,2,opt,name=paramStr,proto3" json:"paramStr,omitempty"`
	SingName             string   `protobuf:"bytes,3,opt,name=sing_name,json=singName,proto3" json:"sing_name,omitempty"`
	TemplateCode         string   `protobuf:"bytes,4,opt,name=template_code,json=templateCode,proto3" json:"template_code,omitempty"`
	KeyId                string   `protobuf:"bytes,5,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	Key                  string   `protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendRequest) Reset()         { *m = SendRequest{} }
func (m *SendRequest) String() string { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()    {}
func (*SendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_604790540bb3bbcc, []int{0}
}

func (m *SendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRequest.Unmarshal(m, b)
}
func (m *SendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRequest.Marshal(b, m, deterministic)
}
func (m *SendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRequest.Merge(m, src)
}
func (m *SendRequest) XXX_Size() int {
	return xxx_messageInfo_SendRequest.Size(m)
}
func (m *SendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendRequest proto.InternalMessageInfo

func (m *SendRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *SendRequest) GetParamStr() string {
	if m != nil {
		return m.ParamStr
	}
	return ""
}

func (m *SendRequest) GetSingName() string {
	if m != nil {
		return m.SingName
	}
	return ""
}

func (m *SendRequest) GetTemplateCode() string {
	if m != nil {
		return m.TemplateCode
	}
	return ""
}

func (m *SendRequest) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

func (m *SendRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type SendResponse struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendResponse) Reset()         { *m = SendResponse{} }
func (m *SendResponse) String() string { return proto.CompactTextString(m) }
func (*SendResponse) ProtoMessage()    {}
func (*SendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_604790540bb3bbcc, []int{1}
}

func (m *SendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendResponse.Unmarshal(m, b)
}
func (m *SendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendResponse.Marshal(b, m, deterministic)
}
func (m *SendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendResponse.Merge(m, src)
}
func (m *SendResponse) XXX_Size() int {
	return xxx_messageInfo_SendResponse.Size(m)
}
func (m *SendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SendResponse proto.InternalMessageInfo

func (m *SendResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*SendRequest)(nil), "app.SendRequest")
	proto.RegisterType((*SendResponse)(nil), "app.SendResponse")
}

func init() { proto.RegisterFile("aliyun.proto", fileDescriptor_604790540bb3bbcc) }

var fileDescriptor_604790540bb3bbcc = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0xa9, 0xb1, 0xa1, 0xbd, 0x46, 0xa9, 0x17, 0x95, 0xa1, 0x6e, 0xa4, 0x82, 0xb8, 0x8a,
	0xa0, 0x4f, 0x20, 0xae, 0xdc, 0xb8, 0x68, 0x1e, 0x20, 0x8c, 0x9d, 0x83, 0x86, 0x64, 0x7e, 0x9c,
	0x99, 0x08, 0x79, 0x25, 0x9f, 0x52, 0x92, 0xa9, 0xd2, 0xdd, 0x9c, 0xef, 0x1c, 0xe6, 0x9e, 0x7b,
	0xa9, 0x90, 0x5d, 0x33, 0xf4, 0xa6, 0x74, 0xde, 0x46, 0xcb, 0x99, 0x74, 0x6e, 0xf3, 0x33, 0xa3,
	0x93, 0x0a, 0x46, 0x6d, 0xf1, 0xd5, 0x23, 0x44, 0xbe, 0xa0, 0xb9, 0xfb, 0xb4, 0x06, 0x62, 0x76,
	0x33, 0xbb, 0x5f, 0x6e, 0x93, 0xe0, 0x35, 0x2d, 0x9c, 0xf4, 0x52, 0x57, 0xd1, 0x8b, 0xa3, 0xc9,
	0xf8, 0xd7, 0x7c, 0x4d, 0xcb, 0xd0, 0x98, 0x8f, 0xda, 0x48, 0x0d, 0x91, 0x25, 0x73, 0x04, 0x6f,
	0x52, 0x83, 0x6f, 0xe9, 0x34, 0x42, 0xbb, 0x4e, 0x46, 0xd4, 0x3b, 0xab, 0x20, 0x8e, 0xa7, 0x40,
	0xf1, 0x07, 0x5f, 0xac, 0x02, 0x5f, 0x52, 0xde, 0x62, 0xa8, 0x1b, 0x25, 0xe6, 0x69, 0x68, 0x8b,
	0xe1, 0x55, 0xf1, 0x8a, 0xb2, 0x16, 0x83, 0xc8, 0x27, 0x36, 0x3e, 0x37, 0x77, 0x54, 0xa4, 0xae,
	0xc1, 0x59, 0x13, 0xc0, 0x57, 0x94, 0x7b, 0x84, 0xbe, 0x8b, 0xfb, 0xb6, 0x7b, 0xf5, 0xf8, 0x4c,
	0x67, 0x63, 0xae, 0xd2, 0xa1, 0x82, 0xff, 0x6e, 0x76, 0xe0, 0x07, 0x5a, 0x04, 0x18, 0x55, 0x07,
	0x1d, 0x78, 0x55, 0x4a, 0xe7, 0xca, 0x83, 0xa5, 0xd7, 0xe7, 0x07, 0x24, 0x7d, 0xfd, 0x9e, 0x4f,
	0x37, 0x7a, 0xfa, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x1b, 0xd8, 0x66, 0x33, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SendSmsServiceClient is the client API for SendSmsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SendSmsServiceClient interface {
	//??????????????????
	SendSms(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error)
}

type sendSmsServiceClient struct {
	cc *grpc.ClientConn
}

func NewSendSmsServiceClient(cc *grpc.ClientConn) SendSmsServiceClient {
	return &sendSmsServiceClient{cc}
}

func (c *sendSmsServiceClient) SendSms(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/app.SendSmsService/send_sms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendSmsServiceServer is the server API for SendSmsService service.
type SendSmsServiceServer interface {
	//??????????????????
	SendSms(context.Context, *SendRequest) (*SendResponse, error)
}

// UnimplementedSendSmsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSendSmsServiceServer struct {
}

func (*UnimplementedSendSmsServiceServer) SendSms(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSms not implemented")
}

func RegisterSendSmsServiceServer(s *grpc.Server, srv SendSmsServiceServer) {
	s.RegisterService(&_SendSmsService_serviceDesc, srv)
}

func _SendSmsService_SendSms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendSmsServiceServer).SendSms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.SendSmsService/SendSms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendSmsServiceServer).SendSms(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SendSmsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.SendSmsService",
	HandlerType: (*SendSmsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "send_sms",
			Handler:    _SendSmsService_SendSms_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aliyun.proto",
}
