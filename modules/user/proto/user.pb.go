// Code generated by protoc-gen-go. DO NOT EDIT.
// source: modules/user/proto/user.proto

package userpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User_Status int32

const (
	User_USER_STATUS_INVALID    User_Status = 0
	User_USER_STATUS_REGISTERED User_Status = 1
	User_USER_STATUS_ACTIVE     User_Status = 2
	User_USER_STATUS_BANNED     User_Status = 3
)

var User_Status_name = map[int32]string{
	0: "USER_STATUS_INVALID",
	1: "USER_STATUS_REGISTERED",
	2: "USER_STATUS_ACTIVE",
	3: "USER_STATUS_BANNED",
}

var User_Status_value = map[string]int32{
	"USER_STATUS_INVALID":    0,
	"USER_STATUS_REGISTERED": 1,
	"USER_STATUS_ACTIVE":     2,
	"USER_STATUS_BANNED":     3,
}

func (x User_Status) String() string {
	return proto.EnumName(User_Status_name, int32(x))
}

func (User_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_234d831259a1c5d0, []int{0, 0}
}

type User struct {
	Id                   int64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username             string      `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email                string      `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Status               User_Status `protobuf:"varint,4,opt,name=status,proto3,enum=user.User_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_234d831259a1c5d0, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetStatus() User_Status {
	if m != nil {
		return m.Status
	}
	return User_USER_STATUS_INVALID
}

type LoginRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_234d831259a1c5d0, []int{1}
}

func (m *LoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginRequest.Unmarshal(m, b)
}
func (m *LoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginRequest.Marshal(b, m, deterministic)
}
func (m *LoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginRequest.Merge(m, src)
}
func (m *LoginRequest) XXX_Size() int {
	return xxx_messageInfo_LoginRequest.Size(m)
}
func (m *LoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LogoutRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutRequest) Reset()         { *m = LogoutRequest{} }
func (m *LogoutRequest) String() string { return proto.CompactTextString(m) }
func (*LogoutRequest) ProtoMessage()    {}
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_234d831259a1c5d0, []int{2}
}

func (m *LogoutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutRequest.Unmarshal(m, b)
}
func (m *LogoutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutRequest.Marshal(b, m, deterministic)
}
func (m *LogoutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutRequest.Merge(m, src)
}
func (m *LogoutRequest) XXX_Size() int {
	return xxx_messageInfo_LogoutRequest.Size(m)
}
func (m *LogoutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutRequest proto.InternalMessageInfo

type NoopResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoopResponse) Reset()         { *m = NoopResponse{} }
func (m *NoopResponse) String() string { return proto.CompactTextString(m) }
func (*NoopResponse) ProtoMessage()    {}
func (*NoopResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_234d831259a1c5d0, []int{3}
}

func (m *NoopResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoopResponse.Unmarshal(m, b)
}
func (m *NoopResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoopResponse.Marshal(b, m, deterministic)
}
func (m *NoopResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoopResponse.Merge(m, src)
}
func (m *NoopResponse) XXX_Size() int {
	return xxx_messageInfo_NoopResponse.Size(m)
}
func (m *NoopResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NoopResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NoopResponse proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("user.User_Status", User_Status_name, User_Status_value)
	proto.RegisterType((*User)(nil), "user.User")
	proto.RegisterType((*LoginRequest)(nil), "user.LoginRequest")
	proto.RegisterType((*LogoutRequest)(nil), "user.LogoutRequest")
	proto.RegisterType((*NoopResponse)(nil), "user.NoopResponse")
}

func init() { proto.RegisterFile("modules/user/proto/user.proto", fileDescriptor_234d831259a1c5d0) }

var fileDescriptor_234d831259a1c5d0 = []byte{
	// 397 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x65, 0x9d, 0xd4, 0x4a, 0x46, 0x21, 0x84, 0x69, 0x55, 0x22, 0x0b, 0x44, 0xb4, 0xa7, 0xc0,
	0x21, 0x16, 0xe5, 0x86, 0xc4, 0x21, 0x21, 0x06, 0x45, 0x8a, 0xa2, 0x6a, 0x9d, 0xf4, 0xc0, 0xa5,
	0x72, 0xf1, 0xca, 0xb2, 0x64, 0x7b, 0x8c, 0x77, 0x0d, 0xe2, 0xca, 0x27, 0xc0, 0xa7, 0xf1, 0x07,
	0x88, 0x0f, 0x41, 0xbb, 0x76, 0x5b, 0x37, 0xb7, 0x99, 0x79, 0x6f, 0xdf, 0x9b, 0x37, 0x36, 0xbc,
	0xc8, 0x29, 0xae, 0x33, 0xa9, 0xfc, 0x5a, 0xc9, 0xca, 0x2f, 0x2b, 0xd2, 0x64, 0xcb, 0x85, 0x2d,
	0xb1, 0x6f, 0x6a, 0xef, 0x79, 0x42, 0x94, 0x64, 0xd2, 0x8f, 0xca, 0xd4, 0x8f, 0x8a, 0x82, 0x74,
	0xa4, 0x53, 0x2a, 0x54, 0xc3, 0xe1, 0x7f, 0x19, 0xf4, 0x0f, 0x4a, 0x56, 0x38, 0x06, 0x27, 0x8d,
	0xa7, 0x6c, 0xc6, 0xe6, 0x3d, 0xe1, 0xa4, 0x31, 0x7a, 0x30, 0x30, 0xcf, 0x8b, 0x28, 0x97, 0x53,
	0x67, 0xc6, 0xe6, 0x43, 0x71, 0xd7, 0xe3, 0x19, 0x9c, 0xc8, 0x3c, 0x4a, 0xb3, 0x69, 0xcf, 0x02,
	0x4d, 0x83, 0xaf, 0xc0, 0x55, 0x3a, 0xd2, 0xb5, 0x9a, 0xf6, 0x67, 0x6c, 0x3e, 0xbe, 0x78, 0xba,
	0xb0, 0xbb, 0x18, 0xf5, 0x45, 0x68, 0x01, 0xd1, 0x12, 0x78, 0x0e, 0x6e, 0x33, 0xc1, 0x67, 0x70,
	0x7a, 0x08, 0x03, 0x71, 0x1d, 0xee, 0x97, 0xfb, 0x43, 0x78, 0xbd, 0xd9, 0x5d, 0x2d, 0xb7, 0x9b,
	0xf5, 0xe4, 0x11, 0x7a, 0x70, 0xde, 0x05, 0x44, 0xf0, 0x69, 0x13, 0xee, 0x03, 0x11, 0xac, 0x27,
	0x0c, 0xcf, 0x01, 0xbb, 0xd8, 0xf2, 0xc3, 0x7e, 0x73, 0x15, 0x4c, 0x9c, 0xe3, 0xf9, 0x6a, 0xb9,
	0xdb, 0x05, 0xeb, 0x49, 0x8f, 0x7f, 0x84, 0xd1, 0x96, 0x92, 0xb4, 0x10, 0xf2, 0x6b, 0x2d, 0x95,
	0x7e, 0x90, 0x8d, 0x1d, 0x65, 0xf3, 0x60, 0x50, 0x46, 0x4a, 0x7d, 0xa7, 0x2a, 0xbe, 0xcd, 0x7d,
	0xdb, 0xf3, 0x27, 0xf0, 0x78, 0x4b, 0x09, 0xd5, 0xba, 0x15, 0xe2, 0x63, 0x18, 0xed, 0x88, 0x4a,
	0x21, 0x55, 0x49, 0x85, 0x92, 0x17, 0xbf, 0x18, 0x80, 0xc9, 0x1b, 0xfe, 0x50, 0x5a, 0xe6, 0xf8,
	0x1e, 0x4e, 0xac, 0x2f, 0x62, 0x73, 0x8a, 0xee, 0x12, 0x1e, 0xdc, 0x9f, 0x87, 0x9f, 0xfd, 0xfc,
	0xf3, 0xef, 0xb7, 0x33, 0xe6, 0x43, 0xff, 0xdb, 0x1b, 0x3f, 0x33, 0xac, 0x77, 0xec, 0x35, 0x06,
	0xe0, 0x36, 0x76, 0x78, 0x7a, 0xf7, 0xfe, 0xde, 0xdc, 0x6b, 0x45, 0xbb, 0x0b, 0x70, 0xb4, 0x42,
	0x23, 0x84, 0x56, 0x88, 0x6a, 0xbd, 0x7a, 0x09, 0x83, 0x2f, 0x94, 0x5b, 0xf2, 0x6a, 0x68, 0xec,
	0x2e, 0xcd, 0x97, 0xbf, 0x64, 0x9f, 0x5d, 0x33, 0x2a, 0x6f, 0x6e, 0x5c, 0xfb, 0x2b, 0xbc, 0xfd,
	0x1f, 0x00, 0x00, 0xff, 0xff, 0x95, 0xbc, 0x92, 0xc1, 0x4f, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserSystemClient is the client API for UserSystem service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserSystemClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*User, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*NoopResponse, error)
}

type userSystemClient struct {
	cc *grpc.ClientConn
}

func NewUserSystemClient(cc *grpc.ClientConn) UserSystemClient {
	return &userSystemClient{cc}
}

func (c *userSystemClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/user.UserSystem/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userSystemClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*NoopResponse, error) {
	out := new(NoopResponse)
	err := c.cc.Invoke(ctx, "/user.UserSystem/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserSystemServer is the server API for UserSystem service.
type UserSystemServer interface {
	Login(context.Context, *LoginRequest) (*User, error)
	Logout(context.Context, *LogoutRequest) (*NoopResponse, error)
}

func RegisterUserSystemServer(s *grpc.Server, srv UserSystemServer) {
	s.RegisterService(&_UserSystem_serviceDesc, srv)
}

func _UserSystem_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSystemServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserSystem/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSystemServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserSystem_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSystemServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserSystem/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSystemServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserSystem_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserSystem",
	HandlerType: (*UserSystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _UserSystem_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _UserSystem_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "modules/user/proto/user.proto",
}
