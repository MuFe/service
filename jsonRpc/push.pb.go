// Code generated by protoc-gen-go. DO NOT EDIT.
// source: push.proto

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

type PushRequest struct {
	DeviceList           []string `protobuf:"bytes,1,rep,name=device_list,json=deviceList,proto3" json:"device_list,omitempty"`
	Content              string   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Message              string   `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushRequest) Reset()         { *m = PushRequest{} }
func (m *PushRequest) String() string { return proto.CompactTextString(m) }
func (*PushRequest) ProtoMessage()    {}
func (*PushRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1e4bfd2e9d102bb, []int{0}
}

func (m *PushRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushRequest.Unmarshal(m, b)
}
func (m *PushRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushRequest.Marshal(b, m, deterministic)
}
func (m *PushRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushRequest.Merge(m, src)
}
func (m *PushRequest) XXX_Size() int {
	return xxx_messageInfo_PushRequest.Size(m)
}
func (m *PushRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PushRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PushRequest proto.InternalMessageInfo

func (m *PushRequest) GetDeviceList() []string {
	if m != nil {
		return m.DeviceList
	}
	return nil
}

func (m *PushRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *PushRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *PushRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type PushResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushResponse) Reset()         { *m = PushResponse{} }
func (m *PushResponse) String() string { return proto.CompactTextString(m) }
func (*PushResponse) ProtoMessage()    {}
func (*PushResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1e4bfd2e9d102bb, []int{1}
}

func (m *PushResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushResponse.Unmarshal(m, b)
}
func (m *PushResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushResponse.Marshal(b, m, deterministic)
}
func (m *PushResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushResponse.Merge(m, src)
}
func (m *PushResponse) XXX_Size() int {
	return xxx_messageInfo_PushResponse.Size(m)
}
func (m *PushResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PushResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PushResponse proto.InternalMessageInfo

func (m *PushResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetPhoneRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPhoneRequest) Reset()         { *m = GetPhoneRequest{} }
func (m *GetPhoneRequest) String() string { return proto.CompactTextString(m) }
func (*GetPhoneRequest) ProtoMessage()    {}
func (*GetPhoneRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1e4bfd2e9d102bb, []int{2}
}

func (m *GetPhoneRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPhoneRequest.Unmarshal(m, b)
}
func (m *GetPhoneRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPhoneRequest.Marshal(b, m, deterministic)
}
func (m *GetPhoneRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPhoneRequest.Merge(m, src)
}
func (m *GetPhoneRequest) XXX_Size() int {
	return xxx_messageInfo_GetPhoneRequest.Size(m)
}
func (m *GetPhoneRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPhoneRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPhoneRequest proto.InternalMessageInfo

func (m *GetPhoneRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type GetPhoneResponse struct {
	Phone                string   `protobuf:"bytes,1,opt,name=phone,proto3" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPhoneResponse) Reset()         { *m = GetPhoneResponse{} }
func (m *GetPhoneResponse) String() string { return proto.CompactTextString(m) }
func (*GetPhoneResponse) ProtoMessage()    {}
func (*GetPhoneResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1e4bfd2e9d102bb, []int{3}
}

func (m *GetPhoneResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPhoneResponse.Unmarshal(m, b)
}
func (m *GetPhoneResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPhoneResponse.Marshal(b, m, deterministic)
}
func (m *GetPhoneResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPhoneResponse.Merge(m, src)
}
func (m *GetPhoneResponse) XXX_Size() int {
	return xxx_messageInfo_GetPhoneResponse.Size(m)
}
func (m *GetPhoneResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPhoneResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPhoneResponse proto.InternalMessageInfo

func (m *GetPhoneResponse) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func init() {
	proto.RegisterType((*PushRequest)(nil), "app.PushRequest")
	proto.RegisterType((*PushResponse)(nil), "app.PushResponse")
	proto.RegisterType((*GetPhoneRequest)(nil), "app.GetPhoneRequest")
	proto.RegisterType((*GetPhoneResponse)(nil), "app.GetPhoneResponse")
}

func init() { proto.RegisterFile("push.proto", fileDescriptor_d1e4bfd2e9d102bb) }

var fileDescriptor_d1e4bfd2e9d102bb = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xbd, 0x4e, 0xc3, 0x40,
	0x10, 0x84, 0x65, 0x4c, 0x40, 0xde, 0x44, 0x22, 0x9c, 0x8c, 0x74, 0x4a, 0x43, 0xe4, 0x06, 0x57,
	0x2e, 0x48, 0xc3, 0x1b, 0xd0, 0x50, 0x44, 0xe6, 0x01, 0x2c, 0x13, 0x56, 0xb1, 0x45, 0xb8, 0x5b,
	0xb2, 0x6b, 0x44, 0xc1, 0xc3, 0xa3, 0xfb, 0xb1, 0x48, 0x5c, 0xce, 0xf8, 0xf3, 0xcc, 0xee, 0x1e,
	0x00, 0x0d, 0xdc, 0x55, 0x74, 0xb4, 0x62, 0x55, 0xda, 0x12, 0x15, 0x3f, 0x30, 0xdf, 0x0e, 0xdc,
	0xd5, 0xf8, 0x35, 0x20, 0x8b, 0xba, 0x87, 0xf9, 0x3b, 0x7e, 0xf7, 0x3b, 0x6c, 0x0e, 0x3d, 0x8b,
	0x4e, 0xd6, 0x69, 0x99, 0xd5, 0x10, 0xac, 0x97, 0x9e, 0x45, 0x69, 0xb8, 0xde, 0x59, 0x23, 0x68,
	0x44, 0x5f, 0xac, 0x93, 0x32, 0xab, 0x47, 0xa9, 0x72, 0x98, 0x49, 0x2f, 0x07, 0xd4, 0xa9, 0xf7,
	0x83, 0x70, 0xfc, 0x27, 0x32, 0xb7, 0x7b, 0xd4, 0x97, 0x81, 0x8f, 0xb2, 0x28, 0x61, 0x11, 0x9a,
	0x99, 0xac, 0xe1, 0x33, 0x32, 0x39, 0x27, 0x1f, 0xe0, 0xe6, 0x19, 0x65, 0xdb, 0x59, 0x83, 0xe3,
	0x9c, 0xae, 0xcc, 0x7e, 0xa0, 0x89, 0x68, 0x10, 0x45, 0x09, 0xcb, 0x7f, 0x30, 0xc6, 0xe6, 0x30,
	0x23, 0x67, 0x8c, 0xa4, 0x17, 0x8f, 0xbf, 0x61, 0xed, 0x57, 0x3c, 0xba, 0xcd, 0xd4, 0x06, 0x16,
	0xee, 0x30, 0x4d, 0x6c, 0x54, 0xcb, 0xaa, 0x25, 0xaa, 0x4e, 0x0e, 0xb3, 0xba, 0x3d, 0x71, 0x62,
	0xf2, 0x13, 0x64, 0x7b, 0x94, 0xc6, 0x07, 0xaa, 0xdc, 0x7f, 0x9f, 0x8c, 0xb9, 0xba, 0x9b, 0xb8,
	0xe1, 0xcf, 0xb7, 0x2b, 0xff, 0x00, 0x9b, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x81, 0x90,
	0xe4, 0x8e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PushServiceClient is the client API for PushService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PushServiceClient interface {
	//??????????????????
	PushMessage(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error)
	GetPhone(ctx context.Context, in *GetPhoneRequest, opts ...grpc.CallOption) (*GetPhoneResponse, error)
}

type pushServiceClient struct {
	cc *grpc.ClientConn
}

func NewPushServiceClient(cc *grpc.ClientConn) PushServiceClient {
	return &pushServiceClient{cc}
}

func (c *pushServiceClient) PushMessage(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error) {
	out := new(PushResponse)
	err := c.cc.Invoke(ctx, "/app.PushService/push_message", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pushServiceClient) GetPhone(ctx context.Context, in *GetPhoneRequest, opts ...grpc.CallOption) (*GetPhoneResponse, error) {
	out := new(GetPhoneResponse)
	err := c.cc.Invoke(ctx, "/app.PushService/get_phone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PushServiceServer is the server API for PushService service.
type PushServiceServer interface {
	//??????????????????
	PushMessage(context.Context, *PushRequest) (*PushResponse, error)
	GetPhone(context.Context, *GetPhoneRequest) (*GetPhoneResponse, error)
}

// UnimplementedPushServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPushServiceServer struct {
}

func (*UnimplementedPushServiceServer) PushMessage(ctx context.Context, req *PushRequest) (*PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessage not implemented")
}
func (*UnimplementedPushServiceServer) GetPhone(ctx context.Context, req *GetPhoneRequest) (*GetPhoneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhone not implemented")
}

func RegisterPushServiceServer(s *grpc.Server, srv PushServiceServer) {
	s.RegisterService(&_PushService_serviceDesc, srv)
}

func _PushService_PushMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PushServiceServer).PushMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.PushService/PushMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PushServiceServer).PushMessage(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PushService_GetPhone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPhoneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PushServiceServer).GetPhone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.PushService/GetPhone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PushServiceServer).GetPhone(ctx, req.(*GetPhoneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PushService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.PushService",
	HandlerType: (*PushServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "push_message",
			Handler:    _PushService_PushMessage_Handler,
		},
		{
			MethodName: "get_phone",
			Handler:    _PushService_GetPhone_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "push.proto",
}
