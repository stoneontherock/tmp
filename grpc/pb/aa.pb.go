// Code generated by protoc-gen-go. DO NOT EDIT.
// source: aa.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type AuthenticateIn struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Action               string   `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Domain               string   `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Resource             string   `protobuf:"bytes,4,opt,name=resource,proto3" json:"resource,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthenticateIn) Reset()         { *m = AuthenticateIn{} }
func (m *AuthenticateIn) String() string { return proto.CompactTextString(m) }
func (*AuthenticateIn) ProtoMessage()    {}
func (*AuthenticateIn) Descriptor() ([]byte, []int) {
	return fileDescriptor_e6ba9c53fd79526f, []int{0}
}

func (m *AuthenticateIn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticateIn.Unmarshal(m, b)
}
func (m *AuthenticateIn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticateIn.Marshal(b, m, deterministic)
}
func (m *AuthenticateIn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticateIn.Merge(m, src)
}
func (m *AuthenticateIn) XXX_Size() int {
	return xxx_messageInfo_AuthenticateIn.Size(m)
}
func (m *AuthenticateIn) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticateIn.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticateIn proto.InternalMessageInfo

func (m *AuthenticateIn) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *AuthenticateIn) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *AuthenticateIn) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *AuthenticateIn) GetResource() string {
	if m != nil {
		return m.Resource
	}
	return ""
}

type AuthenticateOut struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthenticateOut) Reset()         { *m = AuthenticateOut{} }
func (m *AuthenticateOut) String() string { return proto.CompactTextString(m) }
func (*AuthenticateOut) ProtoMessage()    {}
func (*AuthenticateOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_e6ba9c53fd79526f, []int{1}
}

func (m *AuthenticateOut) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticateOut.Unmarshal(m, b)
}
func (m *AuthenticateOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticateOut.Marshal(b, m, deterministic)
}
func (m *AuthenticateOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticateOut.Merge(m, src)
}
func (m *AuthenticateOut) XXX_Size() int {
	return xxx_messageInfo_AuthenticateOut.Size(m)
}
func (m *AuthenticateOut) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticateOut.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticateOut proto.InternalMessageInfo

type AuthorizeIn struct {
	UserName             string   `protobuf:"bytes,1,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizeIn) Reset()         { *m = AuthorizeIn{} }
func (m *AuthorizeIn) String() string { return proto.CompactTextString(m) }
func (*AuthorizeIn) ProtoMessage()    {}
func (*AuthorizeIn) Descriptor() ([]byte, []int) {
	return fileDescriptor_e6ba9c53fd79526f, []int{2}
}

func (m *AuthorizeIn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthorizeIn.Unmarshal(m, b)
}
func (m *AuthorizeIn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthorizeIn.Marshal(b, m, deterministic)
}
func (m *AuthorizeIn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizeIn.Merge(m, src)
}
func (m *AuthorizeIn) XXX_Size() int {
	return xxx_messageInfo_AuthorizeIn.Size(m)
}
func (m *AuthorizeIn) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizeIn.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizeIn proto.InternalMessageInfo

func (m *AuthorizeIn) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *AuthorizeIn) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type AuthorizeOut struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizeOut) Reset()         { *m = AuthorizeOut{} }
func (m *AuthorizeOut) String() string { return proto.CompactTextString(m) }
func (*AuthorizeOut) ProtoMessage()    {}
func (*AuthorizeOut) Descriptor() ([]byte, []int) {
	return fileDescriptor_e6ba9c53fd79526f, []int{3}
}

func (m *AuthorizeOut) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthorizeOut.Unmarshal(m, b)
}
func (m *AuthorizeOut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthorizeOut.Marshal(b, m, deterministic)
}
func (m *AuthorizeOut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizeOut.Merge(m, src)
}
func (m *AuthorizeOut) XXX_Size() int {
	return xxx_messageInfo_AuthorizeOut.Size(m)
}
func (m *AuthorizeOut) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizeOut.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizeOut proto.InternalMessageInfo

func (m *AuthorizeOut) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthenticateIn)(nil), "pb.AuthenticateIn")
	proto.RegisterType((*AuthenticateOut)(nil), "pb.AuthenticateOut")
	proto.RegisterType((*AuthorizeIn)(nil), "pb.AuthorizeIn")
	proto.RegisterType((*AuthorizeOut)(nil), "pb.AuthorizeOut")
}

func init() { proto.RegisterFile("aa.proto", fileDescriptor_e6ba9c53fd79526f) }

var fileDescriptor_e6ba9c53fd79526f = []byte{
	// 236 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xc4, 0x40,
	0x0c, 0x86, 0xdd, 0xaa, 0x4b, 0x1b, 0xc5, 0xd5, 0x28, 0x52, 0xea, 0x45, 0x8a, 0x07, 0x4f, 0x45,
	0xf4, 0xe4, 0xb1, 0x17, 0xc1, 0x8b, 0x0b, 0xbe, 0x80, 0xa4, 0xdd, 0x80, 0x45, 0x3a, 0xa9, 0xd3,
	0x0c, 0x82, 0x4f, 0x2f, 0xd3, 0xce, 0x0e, 0xae, 0x78, 0xfc, 0xbf, 0xc9, 0x24, 0x5f, 0x02, 0x29,
	0x51, 0x35, 0x58, 0x51, 0xc1, 0x64, 0x68, 0x4a, 0x0b, 0x27, 0xb5, 0xd3, 0x77, 0x36, 0xda, 0xb5,
	0xa4, 0xfc, 0x6c, 0xf0, 0x02, 0x0e, 0x55, 0x3e, 0xd8, 0xe4, 0x8b, 0xeb, 0xc5, 0x6d, 0xf6, 0x3a,
	0x07, 0xbc, 0x84, 0x25, 0xb5, 0xda, 0x89, 0xc9, 0x93, 0x09, 0x87, 0xe4, 0xf9, 0x46, 0x7a, 0xea,
	0x4c, 0xbe, 0x3f, 0xf3, 0x39, 0x61, 0x01, 0xa9, 0xe5, 0x51, 0x9c, 0x6d, 0x39, 0x3f, 0x98, 0x5e,
	0x62, 0x2e, 0xcf, 0x60, 0xf5, 0x7b, 0xe6, 0xda, 0x69, 0xf9, 0x04, 0x47, 0x1e, 0x89, 0xed, 0xbe,
	0xbd, 0xc3, 0x15, 0x64, 0x6e, 0x64, 0xfb, 0x66, 0xa8, 0xe7, 0xe0, 0x91, 0x7a, 0xf0, 0x42, 0x3d,
	0xfb, 0xd6, 0x03, 0x8d, 0xe3, 0x97, 0xd8, 0x4d, 0x90, 0x89, 0xb9, 0xbc, 0x81, 0xe3, 0xd8, 0x67,
	0xed, 0xf4, 0xff, 0x65, 0xee, 0x3f, 0x21, 0xa9, 0x6b, 0x7c, 0x9c, 0x6b, 0xb7, 0x1a, 0x88, 0xd5,
	0xd0, 0x54, 0xbb, 0xc7, 0x28, 0xce, 0xff, 0x32, 0x2f, 0xbb, 0x87, 0x77, 0x90, 0xc5, 0x31, 0xb8,
	0xda, 0xd6, 0x04, 0xfb, 0xe2, 0x74, 0x07, 0x4c, 0x3f, 0x9a, 0xe5, 0x74, 0xf2, 0x87, 0x9f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xde, 0xfb, 0xbc, 0x1e, 0x7e, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AAClient is the client API for AA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AAClient interface {
	Authenticate(ctx context.Context, in *AuthenticateIn, opts ...grpc.CallOption) (*AuthenticateOut, error)
	Authorize(ctx context.Context, in *AuthorizeIn, opts ...grpc.CallOption) (*AuthorizeOut, error)
}

type aAClient struct {
	cc *grpc.ClientConn
}

func NewAAClient(cc *grpc.ClientConn) AAClient {
	return &aAClient{cc}
}

func (c *aAClient) Authenticate(ctx context.Context, in *AuthenticateIn, opts ...grpc.CallOption) (*AuthenticateOut, error) {
	out := new(AuthenticateOut)
	err := c.cc.Invoke(ctx, "/pb.AA/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aAClient) Authorize(ctx context.Context, in *AuthorizeIn, opts ...grpc.CallOption) (*AuthorizeOut, error) {
	out := new(AuthorizeOut)
	err := c.cc.Invoke(ctx, "/pb.AA/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AAServer is the server API for AA service.
type AAServer interface {
	Authenticate(context.Context, *AuthenticateIn) (*AuthenticateOut, error)
	Authorize(context.Context, *AuthorizeIn) (*AuthorizeOut, error)
}

func RegisterAAServer(s *grpc.Server, srv AAServer) {
	s.RegisterService(&_AA_serviceDesc, srv)
}

func _AA_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AAServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AA/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AAServer).Authenticate(ctx, req.(*AuthenticateIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _AA_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AAServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AA/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AAServer).Authorize(ctx, req.(*AuthorizeIn))
	}
	return interceptor(ctx, in, info, handler)
}

var _AA_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AA",
	HandlerType: (*AAServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AA_Authenticate_Handler,
		},
		{
			MethodName: "Authorize",
			Handler:    _AA_Authorize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aa.proto",
}