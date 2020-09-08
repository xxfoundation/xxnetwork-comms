// Code generated by protoc-gen-go. DO NOT EDIT.
// source: interconnect.proto

package interconnect

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	messages "gitlab.com/xx_network/comms/messages"
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

// The Network Definition File is defined as a
// JSON structure in primitives/ndf.
type NDF struct {
	Ndf                  []byte   `protobuf:"bytes,1,opt,name=Ndf,proto3" json:"Ndf,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NDF) Reset()         { *m = NDF{} }
func (m *NDF) String() string { return proto.CompactTextString(m) }
func (*NDF) ProtoMessage()    {}
func (*NDF) Descriptor() ([]byte, []int) {
	return fileDescriptor_076c6e9f11b66192, []int{0}
}

func (m *NDF) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NDF.Unmarshal(m, b)
}
func (m *NDF) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NDF.Marshal(b, m, deterministic)
}
func (m *NDF) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NDF.Merge(m, src)
}
func (m *NDF) XXX_Size() int {
	return xxx_messageInfo_NDF.Size(m)
}
func (m *NDF) XXX_DiscardUnknown() {
	xxx_messageInfo_NDF.DiscardUnknown(m)
}

var xxx_messageInfo_NDF proto.InternalMessageInfo

func (m *NDF) GetNdf() []byte {
	if m != nil {
		return m.Ndf
	}
	return nil
}

func init() {
	proto.RegisterType((*NDF)(nil), "interconnect.NDF")
}

func init() { proto.RegisterFile("interconnect.proto", fileDescriptor_076c6e9f11b66192) }

var fileDescriptor_076c6e9f11b66192 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xcc, 0x2b, 0x49,
	0x2d, 0x4a, 0xce, 0xcf, 0xcb, 0x4b, 0x4d, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x41, 0x16, 0x93, 0x32, 0x4e, 0xcf, 0x2c, 0xc9, 0x49, 0x4c, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xaf,
	0xa8, 0x88, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x4f, 0xce, 0xcf, 0xcd, 0x2d, 0xd6,
	0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x45, 0x30, 0x20, 0x46, 0x28, 0x89, 0x73, 0x31, 0xfb,
	0xb9, 0xb8, 0x09, 0x09, 0x70, 0x31, 0xfb, 0xa5, 0xa4, 0x49, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04,
	0x81, 0x98, 0x46, 0xd6, 0x5c, 0x3c, 0x9e, 0x48, 0xa6, 0x0b, 0x69, 0x73, 0xb1, 0xb9, 0xa7, 0x96,
	0x80, 0xd4, 0xf2, 0xe9, 0xc1, 0xcd, 0x08, 0xc8, 0xcc, 0x4b, 0x97, 0x12, 0xd4, 0x43, 0x71, 0x9a,
	0x9f, 0x8b, 0x5b, 0x12, 0x1b, 0xd8, 0x70, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1b, 0x2e,
	0x36, 0xe8, 0xb5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InterconnectClient is the client API for Interconnect service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InterconnectClient interface {
	GetNDF(ctx context.Context, in *messages.Ping, opts ...grpc.CallOption) (*NDF, error)
}

type interconnectClient struct {
	cc *grpc.ClientConn
}

func NewInterconnectClient(cc *grpc.ClientConn) InterconnectClient {
	return &interconnectClient{cc}
}

func (c *interconnectClient) GetNDF(ctx context.Context, in *messages.Ping, opts ...grpc.CallOption) (*NDF, error) {
	out := new(NDF)
	err := c.cc.Invoke(ctx, "/interconnect.Interconnect/GetNDF", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InterconnectServer is the server API for Interconnect service.
type InterconnectServer interface {
	GetNDF(context.Context, *messages.Ping) (*NDF, error)
}

// UnimplementedInterconnectServer can be embedded to have forward compatible implementations.
type UnimplementedInterconnectServer struct {
}

func (*UnimplementedInterconnectServer) GetNDF(ctx context.Context, req *messages.Ping) (*NDF, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNDF not implemented")
}

func RegisterInterconnectServer(s *grpc.Server, srv InterconnectServer) {
	s.RegisterService(&_Interconnect_serviceDesc, srv)
}

func _Interconnect_GetNDF_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(messages.Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InterconnectServer).GetNDF(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interconnect.Interconnect/GetNDF",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InterconnectServer).GetNDF(ctx, req.(*messages.Ping))
	}
	return interceptor(ctx, in, info, handler)
}

var _Interconnect_serviceDesc = grpc.ServiceDesc{
	ServiceName: "interconnect.Interconnect",
	HandlerType: (*InterconnectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNDF",
			Handler:    _Interconnect_GetNDF_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interconnect.proto",
}