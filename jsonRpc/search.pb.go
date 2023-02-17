// Code generated by protoc-gen-go. DO NOT EDIT.
// source: search.proto

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

type SearchRequest struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Uid                  int64    `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchRequest) Reset()         { *m = SearchRequest{} }
func (m *SearchRequest) String() string { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()    {}
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{0}
}

func (m *SearchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchRequest.Unmarshal(m, b)
}
func (m *SearchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchRequest.Marshal(b, m, deterministic)
}
func (m *SearchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchRequest.Merge(m, src)
}
func (m *SearchRequest) XXX_Size() int {
	return xxx_messageInfo_SearchRequest.Size(m)
}
func (m *SearchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SearchRequest proto.InternalMessageInfo

func (m *SearchRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *SearchRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type SearchResponse struct {
	List                 []*SearchData `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	Hot                  []*SearchData `protobuf:"bytes,2,rep,name=hot,proto3" json:"hot,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SearchResponse) Reset()         { *m = SearchResponse{} }
func (m *SearchResponse) String() string { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()    {}
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{1}
}

func (m *SearchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchResponse.Unmarshal(m, b)
}
func (m *SearchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchResponse.Marshal(b, m, deterministic)
}
func (m *SearchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchResponse.Merge(m, src)
}
func (m *SearchResponse) XXX_Size() int {
	return xxx_messageInfo_SearchResponse.Size(m)
}
func (m *SearchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SearchResponse proto.InternalMessageInfo

func (m *SearchResponse) GetList() []*SearchData {
	if m != nil {
		return m.List
	}
	return nil
}

func (m *SearchResponse) GetHot() []*SearchData {
	if m != nil {
		return m.Hot
	}
	return nil
}

type SearchData struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TodayNumber          int64    `protobuf:"varint,2,opt,name=today_number,json=todayNumber,proto3" json:"today_number,omitempty"`
	Number               int64    `protobuf:"varint,3,opt,name=number,proto3" json:"number,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchData) Reset()         { *m = SearchData{} }
func (m *SearchData) String() string { return proto.CompactTextString(m) }
func (*SearchData) ProtoMessage()    {}
func (*SearchData) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{2}
}

func (m *SearchData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchData.Unmarshal(m, b)
}
func (m *SearchData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchData.Marshal(b, m, deterministic)
}
func (m *SearchData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchData.Merge(m, src)
}
func (m *SearchData) XXX_Size() int {
	return xxx_messageInfo_SearchData.Size(m)
}
func (m *SearchData) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchData.DiscardUnknown(m)
}

var xxx_messageInfo_SearchData proto.InternalMessageInfo

func (m *SearchData) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SearchData) GetTodayNumber() int64 {
	if m != nil {
		return m.TodayNumber
	}
	return 0
}

func (m *SearchData) GetNumber() int64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *SearchData) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*SearchRequest)(nil), "app.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "app.SearchResponse")
	proto.RegisterType((*SearchData)(nil), "app.SearchData")
}

func init() { proto.RegisterFile("search.proto", fileDescriptor_453745cff914010e) }

var fileDescriptor_453745cff914010e = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x6f, 0xc2, 0x30,
	0x10, 0x85, 0x95, 0x18, 0x51, 0x71, 0x50, 0x8a, 0xae, 0x52, 0x15, 0x31, 0x41, 0xba, 0x64, 0xca,
	0x00, 0xea, 0xc4, 0xda, 0x4a, 0x9d, 0x3a, 0x84, 0xa5, 0x5b, 0x65, 0x62, 0xab, 0x58, 0x2a, 0xb1,
	0x89, 0x2f, 0x95, 0xf2, 0xcf, 0xfa, 0xf3, 0xaa, 0xd8, 0x4e, 0x0b, 0x12, 0x03, 0x9b, 0xef, 0x7b,
	0xef, 0x49, 0x77, 0xcf, 0x30, 0xb1, 0x92, 0xd7, 0xe5, 0x3e, 0x37, 0xb5, 0x26, 0x8d, 0x8c, 0x1b,
	0x33, 0x07, 0xc1, 0x89, 0x7b, 0x90, 0x6e, 0xe0, 0x76, 0xeb, 0x0c, 0x85, 0x3c, 0x36, 0xd2, 0x12,
	0x26, 0x70, 0x53, 0xea, 0x8a, 0x64, 0x45, 0x49, 0xb4, 0x88, 0xb2, 0x51, 0xd1, 0x8f, 0x38, 0x03,
	0xd6, 0x28, 0x91, 0xc4, 0x8b, 0x28, 0x63, 0x45, 0xf7, 0x4c, 0xdf, 0x61, 0xda, 0x87, 0xad, 0xd1,
	0x95, 0x95, 0xf8, 0x08, 0x83, 0x2f, 0x65, 0xbb, 0x28, 0xcb, 0xc6, 0xab, 0xbb, 0x9c, 0x1b, 0x93,
	0x7b, 0xcb, 0x33, 0x27, 0x5e, 0x38, 0x11, 0x97, 0xc0, 0xf6, 0x9a, 0x92, 0xf8, 0xb2, 0xa7, 0xd3,
	0xd2, 0x23, 0xc0, 0x3f, 0xc2, 0x29, 0xc4, 0x4a, 0xb8, 0x75, 0x58, 0x11, 0x2b, 0x81, 0x4b, 0x98,
	0x90, 0x16, 0xbc, 0xfd, 0xa8, 0x9a, 0xc3, 0x4e, 0xd6, 0x61, 0xa5, 0xb1, 0x63, 0x6f, 0x0e, 0xe1,
	0x03, 0x0c, 0x83, 0xc8, 0x9c, 0x18, 0xa6, 0xd3, 0xf3, 0x06, 0x67, 0xe7, 0xad, 0x7e, 0xa2, 0xbe,
	0x8a, 0xad, 0xac, 0xbf, 0x55, 0x29, 0x71, 0x0d, 0x23, 0x2e, 0x84, 0x67, 0x88, 0x27, 0x7b, 0x86,
	0xae, 0xe6, 0x9e, 0xbd, 0x1c, 0x0c, 0xb5, 0x7f, 0x0d, 0x6c, 0x60, 0xf6, 0x29, 0xc9, 0xfb, 0x5e,
	0x95, 0x25, 0x5d, 0xb7, 0x17, 0xb3, 0xf7, 0x67, 0x2c, 0x84, 0x9f, 0x00, 0x6c, 0x48, 0x56, 0x74,
	0x75, 0x6c, 0x37, 0x74, 0x7f, 0xb9, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x71, 0xf8, 0xe2,
	0xec, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SearchServiceClient is the client API for SearchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SearchServiceClient interface {
	//检查通道连接
	AddSearch(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	GetSearchHistory(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	SearchHint(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type searchServiceClient struct {
	cc *grpc.ClientConn
}

func NewSearchServiceClient(cc *grpc.ClientConn) SearchServiceClient {
	return &searchServiceClient{cc}
}

func (c *searchServiceClient) AddSearch(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.SearchService/addSearch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) GetSearchHistory(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/app.SearchService/getSearchHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchServiceClient) SearchHint(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/app.SearchService/searchHint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServiceServer is the server API for SearchService service.
type SearchServiceServer interface {
	//检查通道连接
	AddSearch(context.Context, *SearchRequest) (*EmptyResponse, error)
	GetSearchHistory(context.Context, *SearchRequest) (*SearchResponse, error)
	SearchHint(context.Context, *SearchRequest) (*SearchResponse, error)
}

// UnimplementedSearchServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSearchServiceServer struct {
}

func (*UnimplementedSearchServiceServer) AddSearch(ctx context.Context, req *SearchRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSearch not implemented")
}
func (*UnimplementedSearchServiceServer) GetSearchHistory(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSearchHistory not implemented")
}
func (*UnimplementedSearchServiceServer) SearchHint(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchHint not implemented")
}

func RegisterSearchServiceServer(s *grpc.Server, srv SearchServiceServer) {
	s.RegisterService(&_SearchService_serviceDesc, srv)
}

func _SearchService_AddSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).AddSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.SearchService/AddSearch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).AddSearch(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_GetSearchHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).GetSearchHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.SearchService/GetSearchHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).GetSearchHistory(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchService_SearchHint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServiceServer).SearchHint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.SearchService/SearchHint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServiceServer).SearchHint(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SearchService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.SearchService",
	HandlerType: (*SearchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "addSearch",
			Handler:    _SearchService_AddSearch_Handler,
		},
		{
			MethodName: "getSearchHistory",
			Handler:    _SearchService_GetSearchHistory_Handler,
		},
		{
			MethodName: "searchHint",
			Handler:    _SearchService_SearchHint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search.proto",
}