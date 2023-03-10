// Code generated by protoc-gen-go. DO NOT EDIT.
// source: web.proto

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

type GetWebInfoRequest struct {
	List                 []int64  `protobuf:"varint,1,rep,packed,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetWebInfoRequest) Reset()         { *m = GetWebInfoRequest{} }
func (m *GetWebInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetWebInfoRequest) ProtoMessage()    {}
func (*GetWebInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{0}
}

func (m *GetWebInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWebInfoRequest.Unmarshal(m, b)
}
func (m *GetWebInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWebInfoRequest.Marshal(b, m, deterministic)
}
func (m *GetWebInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWebInfoRequest.Merge(m, src)
}
func (m *GetWebInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetWebInfoRequest.Size(m)
}
func (m *GetWebInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWebInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetWebInfoRequest proto.InternalMessageInfo

func (m *GetWebInfoRequest) GetList() []int64 {
	if m != nil {
		return m.List
	}
	return nil
}

type EditNewResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditNewResponse) Reset()         { *m = EditNewResponse{} }
func (m *EditNewResponse) String() string { return proto.CompactTextString(m) }
func (*EditNewResponse) ProtoMessage()    {}
func (*EditNewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{1}
}

func (m *EditNewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditNewResponse.Unmarshal(m, b)
}
func (m *EditNewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditNewResponse.Marshal(b, m, deterministic)
}
func (m *EditNewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditNewResponse.Merge(m, src)
}
func (m *EditNewResponse) XXX_Size() int {
	return xxx_messageInfo_EditNewResponse.Size(m)
}
func (m *EditNewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EditNewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EditNewResponse proto.InternalMessageInfo

func (m *EditNewResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetWebInfoResponse struct {
	Content              map[int64]string `protobuf:"bytes,1,rep,name=content,proto3" json:"content,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetWebInfoResponse) Reset()         { *m = GetWebInfoResponse{} }
func (m *GetWebInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetWebInfoResponse) ProtoMessage()    {}
func (*GetWebInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{2}
}

func (m *GetWebInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetWebInfoResponse.Unmarshal(m, b)
}
func (m *GetWebInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetWebInfoResponse.Marshal(b, m, deterministic)
}
func (m *GetWebInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetWebInfoResponse.Merge(m, src)
}
func (m *GetWebInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetWebInfoResponse.Size(m)
}
func (m *GetWebInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetWebInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetWebInfoResponse proto.InternalMessageInfo

func (m *GetWebInfoResponse) GetContent() map[int64]string {
	if m != nil {
		return m.Content
	}
	return nil
}

type ContactUsRequest struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Email                string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Id                   int64    `protobuf:"varint,5,opt,name=id,proto3" json:"id,omitempty"`
	Status               int64    `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContactUsRequest) Reset()         { *m = ContactUsRequest{} }
func (m *ContactUsRequest) String() string { return proto.CompactTextString(m) }
func (*ContactUsRequest) ProtoMessage()    {}
func (*ContactUsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{3}
}

func (m *ContactUsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContactUsRequest.Unmarshal(m, b)
}
func (m *ContactUsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContactUsRequest.Marshal(b, m, deterministic)
}
func (m *ContactUsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContactUsRequest.Merge(m, src)
}
func (m *ContactUsRequest) XXX_Size() int {
	return xxx_messageInfo_ContactUsRequest.Size(m)
}
func (m *ContactUsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ContactUsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ContactUsRequest proto.InternalMessageInfo

func (m *ContactUsRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *ContactUsRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *ContactUsRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ContactUsRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ContactUsRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ContactUsRequest) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

type GetNewsRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Page                 int64    `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Size                 int64    `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	Status               int64    `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	Type                 int64    `protobuf:"varint,5,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetNewsRequest) Reset()         { *m = GetNewsRequest{} }
func (m *GetNewsRequest) String() string { return proto.CompactTextString(m) }
func (*GetNewsRequest) ProtoMessage()    {}
func (*GetNewsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{4}
}

func (m *GetNewsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNewsRequest.Unmarshal(m, b)
}
func (m *GetNewsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNewsRequest.Marshal(b, m, deterministic)
}
func (m *GetNewsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNewsRequest.Merge(m, src)
}
func (m *GetNewsRequest) XXX_Size() int {
	return xxx_messageInfo_GetNewsRequest.Size(m)
}
func (m *GetNewsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNewsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetNewsRequest proto.InternalMessageInfo

func (m *GetNewsRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetNewsRequest) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *GetNewsRequest) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *GetNewsRequest) GetStatus() int64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *GetNewsRequest) GetType() int64 {
	if m != nil {
		return m.Type
	}
	return 0
}

type GetNewsResponse struct {
	List                 []*NewsData `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetNewsResponse) Reset()         { *m = GetNewsResponse{} }
func (m *GetNewsResponse) String() string { return proto.CompactTextString(m) }
func (*GetNewsResponse) ProtoMessage()    {}
func (*GetNewsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{5}
}

func (m *GetNewsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNewsResponse.Unmarshal(m, b)
}
func (m *GetNewsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNewsResponse.Marshal(b, m, deterministic)
}
func (m *GetNewsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNewsResponse.Merge(m, src)
}
func (m *GetNewsResponse) XXX_Size() int {
	return xxx_messageInfo_GetNewsResponse.Size(m)
}
func (m *GetNewsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNewsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetNewsResponse proto.InternalMessageInfo

func (m *GetNewsResponse) GetList() []*NewsData {
	if m != nil {
		return m.List
	}
	return nil
}

type NewsData struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Time                 int64    `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Cover                string   `protobuf:"bytes,4,opt,name=cover,proto3" json:"cover,omitempty"`
	Content              string   `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Source               string   `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
	Type                 int64    `protobuf:"varint,7,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewsData) Reset()         { *m = NewsData{} }
func (m *NewsData) String() string { return proto.CompactTextString(m) }
func (*NewsData) ProtoMessage()    {}
func (*NewsData) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{6}
}

func (m *NewsData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewsData.Unmarshal(m, b)
}
func (m *NewsData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewsData.Marshal(b, m, deterministic)
}
func (m *NewsData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewsData.Merge(m, src)
}
func (m *NewsData) XXX_Size() int {
	return xxx_messageInfo_NewsData.Size(m)
}
func (m *NewsData) XXX_DiscardUnknown() {
	xxx_messageInfo_NewsData.DiscardUnknown(m)
}

var xxx_messageInfo_NewsData proto.InternalMessageInfo

func (m *NewsData) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NewsData) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *NewsData) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *NewsData) GetCover() string {
	if m != nil {
		return m.Cover
	}
	return ""
}

func (m *NewsData) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *NewsData) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *NewsData) GetType() int64 {
	if m != nil {
		return m.Type
	}
	return 0
}

type EditNewTypeRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 int64    `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditNewTypeRequest) Reset()         { *m = EditNewTypeRequest{} }
func (m *EditNewTypeRequest) String() string { return proto.CompactTextString(m) }
func (*EditNewTypeRequest) ProtoMessage()    {}
func (*EditNewTypeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{7}
}

func (m *EditNewTypeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditNewTypeRequest.Unmarshal(m, b)
}
func (m *EditNewTypeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditNewTypeRequest.Marshal(b, m, deterministic)
}
func (m *EditNewTypeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditNewTypeRequest.Merge(m, src)
}
func (m *EditNewTypeRequest) XXX_Size() int {
	return xxx_messageInfo_EditNewTypeRequest.Size(m)
}
func (m *EditNewTypeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditNewTypeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditNewTypeRequest proto.InternalMessageInfo

func (m *EditNewTypeRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EditNewTypeRequest) GetType() int64 {
	if m != nil {
		return m.Type
	}
	return 0
}

type EditNewRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Source               string   `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditNewRequest) Reset()         { *m = EditNewRequest{} }
func (m *EditNewRequest) String() string { return proto.CompactTextString(m) }
func (*EditNewRequest) ProtoMessage()    {}
func (*EditNewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{8}
}

func (m *EditNewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditNewRequest.Unmarshal(m, b)
}
func (m *EditNewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditNewRequest.Marshal(b, m, deterministic)
}
func (m *EditNewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditNewRequest.Merge(m, src)
}
func (m *EditNewRequest) XXX_Size() int {
	return xxx_messageInfo_EditNewRequest.Size(m)
}
func (m *EditNewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditNewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditNewRequest proto.InternalMessageInfo

func (m *EditNewRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EditNewRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *EditNewRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *EditNewRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type EditNewCoverRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cover                string   `protobuf:"bytes,2,opt,name=cover,proto3" json:"cover,omitempty"`
	Prefix               string   `protobuf:"bytes,3,opt,name=prefix,proto3" json:"prefix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditNewCoverRequest) Reset()         { *m = EditNewCoverRequest{} }
func (m *EditNewCoverRequest) String() string { return proto.CompactTextString(m) }
func (*EditNewCoverRequest) ProtoMessage()    {}
func (*EditNewCoverRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_461bb3ac99194e85, []int{9}
}

func (m *EditNewCoverRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditNewCoverRequest.Unmarshal(m, b)
}
func (m *EditNewCoverRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditNewCoverRequest.Marshal(b, m, deterministic)
}
func (m *EditNewCoverRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditNewCoverRequest.Merge(m, src)
}
func (m *EditNewCoverRequest) XXX_Size() int {
	return xxx_messageInfo_EditNewCoverRequest.Size(m)
}
func (m *EditNewCoverRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditNewCoverRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditNewCoverRequest proto.InternalMessageInfo

func (m *EditNewCoverRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EditNewCoverRequest) GetCover() string {
	if m != nil {
		return m.Cover
	}
	return ""
}

func (m *EditNewCoverRequest) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func init() {
	proto.RegisterType((*GetWebInfoRequest)(nil), "app.GetWebInfoRequest")
	proto.RegisterType((*EditNewResponse)(nil), "app.EditNewResponse")
	proto.RegisterType((*GetWebInfoResponse)(nil), "app.GetWebInfoResponse")
	proto.RegisterMapType((map[int64]string)(nil), "app.GetWebInfoResponse.ContentEntry")
	proto.RegisterType((*ContactUsRequest)(nil), "app.ContactUsRequest")
	proto.RegisterType((*GetNewsRequest)(nil), "app.GetNewsRequest")
	proto.RegisterType((*GetNewsResponse)(nil), "app.GetNewsResponse")
	proto.RegisterType((*NewsData)(nil), "app.NewsData")
	proto.RegisterType((*EditNewTypeRequest)(nil), "app.EditNewTypeRequest")
	proto.RegisterType((*EditNewRequest)(nil), "app.EditNewRequest")
	proto.RegisterType((*EditNewCoverRequest)(nil), "app.EditNewCoverRequest")
}

func init() { proto.RegisterFile("web.proto", fileDescriptor_461bb3ac99194e85) }

var fileDescriptor_461bb3ac99194e85 = []byte{
	// 593 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x96, 0xed, 0xfc, 0x90, 0x09, 0xa4, 0x65, 0x1b, 0x82, 0x95, 0x53, 0x6a, 0x21, 0x91, 0x53,
	0x0e, 0xa1, 0xaa, 0x4a, 0x24, 0xb8, 0x84, 0xa8, 0xe2, 0xc2, 0xc1, 0x05, 0xf5, 0xbc, 0x71, 0xa6,
	0xa9, 0x45, 0x62, 0x2f, 0xf6, 0x26, 0x25, 0x5c, 0x79, 0x01, 0x0e, 0x3c, 0x00, 0x0f, 0xc6, 0xc3,
	0xa0, 0xfd, 0x73, 0xd7, 0x89, 0x69, 0x6f, 0x33, 0xe3, 0xf9, 0xf6, 0xfb, 0xe6, 0x9b, 0x5d, 0x43,
	0xeb, 0x0e, 0xe7, 0x23, 0x96, 0xa5, 0x3c, 0x25, 0x1e, 0x65, 0xac, 0x0f, 0x0b, 0xca, 0xa9, 0x2a,
	0x04, 0xaf, 0xe1, 0xf9, 0x25, 0xf2, 0x6b, 0x9c, 0x7f, 0x4c, 0x6e, 0xd2, 0x10, 0xbf, 0x6d, 0x30,
	0xe7, 0x84, 0x40, 0x6d, 0x15, 0xe7, 0xdc, 0x77, 0x06, 0xde, 0xd0, 0x0b, 0x65, 0x1c, 0x9c, 0xc2,
	0xd1, 0x6c, 0x11, 0xf3, 0x4f, 0x78, 0x17, 0x62, 0xce, 0xd2, 0x24, 0x47, 0xd2, 0x01, 0x37, 0x5e,
	0xf8, 0xce, 0xc0, 0x19, 0x7a, 0xa1, 0x1b, 0x2f, 0x82, 0x5f, 0x0e, 0x10, 0xfb, 0x30, 0xdd, 0xf6,
	0x1e, 0x9a, 0x51, 0x9a, 0x70, 0x4c, 0xd4, 0x81, 0xed, 0xf1, 0xab, 0x11, 0x65, 0x6c, 0x74, 0xd8,
	0x39, 0x9a, 0xaa, 0xb6, 0x59, 0xc2, 0xb3, 0x5d, 0x68, 0x40, 0xfd, 0x09, 0x3c, 0xb5, 0x3f, 0x90,
	0x63, 0xf0, 0xbe, 0xe2, 0x4e, 0xf3, 0x8a, 0x90, 0x74, 0xa1, 0xbe, 0xa5, 0xab, 0x0d, 0xfa, 0xee,
	0xc0, 0x19, 0xb6, 0x42, 0x95, 0x4c, 0xdc, 0x0b, 0x27, 0xf8, 0xed, 0xc0, 0xb1, 0x00, 0xd3, 0x88,
	0x7f, 0xc9, 0xcd, 0x78, 0xbe, 0x2d, 0x48, 0x00, 0x4c, 0x2a, 0x0e, 0x62, 0xb7, 0x69, 0x52, 0x1c,
	0x24, 0x13, 0x51, 0xc5, 0x35, 0x8d, 0x57, 0xbe, 0xa7, 0xaa, 0x32, 0x11, 0x26, 0x25, 0x74, 0x8d,
	0x7e, 0x4d, 0x16, 0x65, 0xac, 0x1d, 0xa9, 0x1b, 0x47, 0x48, 0x0f, 0x1a, 0x39, 0xa7, 0x7c, 0x93,
	0xfb, 0x0d, 0x59, 0xd3, 0x59, 0xc0, 0xa1, 0x73, 0x89, 0xc2, 0xcb, 0x42, 0xd3, 0x9e, 0x97, 0xe2,
	0x74, 0x46, 0x97, 0x4a, 0x88, 0x17, 0xca, 0x58, 0xd4, 0xf2, 0xf8, 0x07, 0x4a, 0x19, 0x5e, 0x28,
	0x63, 0x8b, 0xa1, 0x66, 0x33, 0x88, 0x5e, 0xbe, 0x63, 0xa8, 0xb5, 0xc8, 0x38, 0x38, 0x83, 0xa3,
	0x82, 0x55, 0xef, 0xe6, 0xd4, 0xda, 0x74, 0x7b, 0xfc, 0x4c, 0x2e, 0x46, 0x34, 0x7c, 0xa0, 0x9c,
	0xea, 0xc5, 0xff, 0x71, 0xe0, 0x89, 0x29, 0x1d, 0xc8, 0xec, 0x42, 0x9d, 0xc7, 0x7c, 0x55, 0x18,
	0x26, 0x13, 0x49, 0x1e, 0xaf, 0x0b, 0xa1, 0x22, 0x16, 0x9d, 0x51, 0xba, 0xc5, 0x4c, 0xfb, 0xa5,
	0x12, 0x7b, 0x15, 0xf5, 0xf2, 0x2a, 0xc4, 0x60, 0xe9, 0x26, 0x8b, 0x50, 0x5a, 0xd7, 0x0a, 0x75,
	0x56, 0x0c, 0xd6, 0xb4, 0x06, 0xbb, 0x00, 0xa2, 0xef, 0xe6, 0xe7, 0x1d, 0xc3, 0x07, 0x2c, 0x95,
	0x48, 0xd7, 0x42, 0xde, 0x42, 0xa7, 0xb8, 0xd5, 0xd5, 0xa8, 0xea, 0x09, 0xef, 0xd5, 0x79, 0x25,
	0x75, 0xd6, 0x3c, 0xb5, 0xd2, 0x3c, 0xc1, 0x15, 0x9c, 0x68, 0xa6, 0xa9, 0x98, 0xfc, 0x01, 0x3a,
	0x65, 0x93, 0x6b, 0xdb, 0xd4, 0x83, 0x06, 0xcb, 0xf0, 0x26, 0xfe, 0x6e, 0xe8, 0x54, 0x36, 0xfe,
	0xe9, 0x00, 0x5c, 0xe3, 0xfc, 0x0a, 0xb3, 0x6d, 0x1c, 0x21, 0x79, 0x07, 0xb0, 0x2c, 0x5e, 0x15,
	0xe9, 0x1d, 0x3c, 0x33, 0x49, 0xd9, 0x7f, 0xf9, 0x9f, 0xe7, 0x47, 0xce, 0xa1, 0x15, 0x99, 0xb7,
	0x42, 0x5e, 0xc8, 0xae, 0xfd, 0xb7, 0xd3, 0x27, 0xb2, 0x3c, 0x5b, 0x33, 0xbe, 0x33, 0xb8, 0xf1,
	0x5f, 0x17, 0xda, 0xe2, 0x86, 0x18, 0x19, 0x67, 0xd0, 0x5c, 0xaa, 0x7b, 0x46, 0x4e, 0x0c, 0x97,
	0x75, 0xd7, 0xfb, 0xdd, 0x72, 0x51, 0xb3, 0x4f, 0xa0, 0x8d, 0xf7, 0x4b, 0x24, 0x4a, 0xe5, 0xe1,
	0x5a, 0xab, 0x14, 0x08, 0x46, 0x8d, 0xd5, 0x8c, 0xe5, 0xa5, 0x6a, 0xc6, 0xfd, 0xff, 0xd7, 0x5b,
	0x68, 0x09, 0xd4, 0x54, 0xdd, 0x44, 0xbb, 0xc5, 0x5e, 0x51, 0x25, 0xe1, 0xb9, 0x12, 0xab, 0xff,
	0x4b, 0xd5, 0xa4, 0x55, 0xb8, 0x31, 0x34, 0x17, 0xb8, 0xb2, 0xac, 0x79, 0x1c, 0x33, 0x6f, 0xc8,
	0x3f, 0xf5, 0x9b, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x39, 0x4c, 0x7e, 0x7c, 0xc7, 0x05, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WebServiceClient is the client API for WebService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServiceClient interface {
	//??????????????????
	GetWebInfo(ctx context.Context, in *GetWebInfoRequest, opts ...grpc.CallOption) (*GetWebInfoResponse, error)
	ContactUs(ctx context.Context, in *ContactUsRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type webServiceClient struct {
	cc *grpc.ClientConn
}

func NewWebServiceClient(cc *grpc.ClientConn) WebServiceClient {
	return &webServiceClient{cc}
}

func (c *webServiceClient) GetWebInfo(ctx context.Context, in *GetWebInfoRequest, opts ...grpc.CallOption) (*GetWebInfoResponse, error) {
	out := new(GetWebInfoResponse)
	err := c.cc.Invoke(ctx, "/app.WebService/getWebInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webServiceClient) ContactUs(ctx context.Context, in *ContactUsRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.WebService/contactUs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WebServiceServer is the server API for WebService service.
type WebServiceServer interface {
	//??????????????????
	GetWebInfo(context.Context, *GetWebInfoRequest) (*GetWebInfoResponse, error)
	ContactUs(context.Context, *ContactUsRequest) (*EmptyResponse, error)
}

// UnimplementedWebServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWebServiceServer struct {
}

func (*UnimplementedWebServiceServer) GetWebInfo(ctx context.Context, req *GetWebInfoRequest) (*GetWebInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWebInfo not implemented")
}
func (*UnimplementedWebServiceServer) ContactUs(ctx context.Context, req *ContactUsRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ContactUs not implemented")
}

func RegisterWebServiceServer(s *grpc.Server, srv WebServiceServer) {
	s.RegisterService(&_WebService_serviceDesc, srv)
}

func _WebService_GetWebInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWebInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).GetWebInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.WebService/GetWebInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).GetWebInfo(ctx, req.(*GetWebInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebService_ContactUs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContactUsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).ContactUs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.WebService/ContactUs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).ContactUs(ctx, req.(*ContactUsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WebService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.WebService",
	HandlerType: (*WebServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getWebInfo",
			Handler:    _WebService_GetWebInfo_Handler,
		},
		{
			MethodName: "contactUs",
			Handler:    _WebService_ContactUs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "web.proto",
}

// NewsServiceClient is the client API for NewsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NewsServiceClient interface {
	GetNews(ctx context.Context, in *GetNewsRequest, opts ...grpc.CallOption) (*GetNewsResponse, error)
	EditNewType(ctx context.Context, in *EditNewTypeRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	EditNew(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EditNewResponse, error)
	EditCover(ctx context.Context, in *EditNewCoverRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	EditContent(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	DelNews(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type newsServiceClient struct {
	cc *grpc.ClientConn
}

func NewNewsServiceClient(cc *grpc.ClientConn) NewsServiceClient {
	return &newsServiceClient{cc}
}

func (c *newsServiceClient) GetNews(ctx context.Context, in *GetNewsRequest, opts ...grpc.CallOption) (*GetNewsResponse, error) {
	out := new(GetNewsResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/getNews", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) EditNewType(ctx context.Context, in *EditNewTypeRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/editNewType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) EditNew(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EditNewResponse, error) {
	out := new(EditNewResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/editNew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) EditCover(ctx context.Context, in *EditNewCoverRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/editCover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) EditContent(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/editContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *newsServiceClient) DelNews(ctx context.Context, in *EditNewRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/app.NewsService/delNews", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NewsServiceServer is the server API for NewsService service.
type NewsServiceServer interface {
	GetNews(context.Context, *GetNewsRequest) (*GetNewsResponse, error)
	EditNewType(context.Context, *EditNewTypeRequest) (*EmptyResponse, error)
	EditNew(context.Context, *EditNewRequest) (*EditNewResponse, error)
	EditCover(context.Context, *EditNewCoverRequest) (*EmptyResponse, error)
	EditContent(context.Context, *EditNewRequest) (*EmptyResponse, error)
	DelNews(context.Context, *EditNewRequest) (*EmptyResponse, error)
}

// UnimplementedNewsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNewsServiceServer struct {
}

func (*UnimplementedNewsServiceServer) GetNews(ctx context.Context, req *GetNewsRequest) (*GetNewsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNews not implemented")
}
func (*UnimplementedNewsServiceServer) EditNewType(ctx context.Context, req *EditNewTypeRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditNewType not implemented")
}
func (*UnimplementedNewsServiceServer) EditNew(ctx context.Context, req *EditNewRequest) (*EditNewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditNew not implemented")
}
func (*UnimplementedNewsServiceServer) EditCover(ctx context.Context, req *EditNewCoverRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditCover not implemented")
}
func (*UnimplementedNewsServiceServer) EditContent(ctx context.Context, req *EditNewRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditContent not implemented")
}
func (*UnimplementedNewsServiceServer) DelNews(ctx context.Context, req *EditNewRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelNews not implemented")
}

func RegisterNewsServiceServer(s *grpc.Server, srv NewsServiceServer) {
	s.RegisterService(&_NewsService_serviceDesc, srv)
}

func _NewsService_GetNews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).GetNews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/GetNews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).GetNews(ctx, req.(*GetNewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_EditNewType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditNewTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).EditNewType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/EditNewType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).EditNewType(ctx, req.(*EditNewTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_EditNew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditNewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).EditNew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/EditNew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).EditNew(ctx, req.(*EditNewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_EditCover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditNewCoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).EditCover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/EditCover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).EditCover(ctx, req.(*EditNewCoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_EditContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditNewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).EditContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/EditContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).EditContent(ctx, req.(*EditNewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NewsService_DelNews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditNewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewsServiceServer).DelNews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.NewsService/DelNews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewsServiceServer).DelNews(ctx, req.(*EditNewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NewsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "app.NewsService",
	HandlerType: (*NewsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getNews",
			Handler:    _NewsService_GetNews_Handler,
		},
		{
			MethodName: "editNewType",
			Handler:    _NewsService_EditNewType_Handler,
		},
		{
			MethodName: "editNew",
			Handler:    _NewsService_EditNew_Handler,
		},
		{
			MethodName: "editCover",
			Handler:    _NewsService_EditCover_Handler,
		},
		{
			MethodName: "editContent",
			Handler:    _NewsService_EditContent_Handler,
		},
		{
			MethodName: "delNews",
			Handler:    _NewsService_DelNews_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "web.proto",
}
