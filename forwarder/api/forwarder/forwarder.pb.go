// Code generated by protoc-gen-go. DO NOT EDIT.
// source: forwarder.proto

package forwarder

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	crossconnect "github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	connection1 "github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	connection "github.com/networkservicemesh/networkservicemesh/controlplane/api/remote/connection"
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

// Message sent by forwarder module informing NSM of any changes in its
// operations parameters or constraints
type MechanismUpdate struct {
	RemoteMechanisms     []*connection.Mechanism  `protobuf:"bytes,1,rep,name=remote_mechanisms,json=remoteMechanisms,proto3" json:"remote_mechanisms,omitempty"`
	LocalMechanisms      []*connection1.Mechanism `protobuf:"bytes,2,rep,name=local_mechanisms,json=localMechanisms,proto3" json:"local_mechanisms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *MechanismUpdate) Reset()         { *m = MechanismUpdate{} }
func (m *MechanismUpdate) String() string { return proto.CompactTextString(m) }
func (*MechanismUpdate) ProtoMessage()    {}
func (*MechanismUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_19bff53f4d11db23, []int{0}
}

func (m *MechanismUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MechanismUpdate.Unmarshal(m, b)
}
func (m *MechanismUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MechanismUpdate.Marshal(b, m, deterministic)
}
func (m *MechanismUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MechanismUpdate.Merge(m, src)
}
func (m *MechanismUpdate) XXX_Size() int {
	return xxx_messageInfo_MechanismUpdate.Size(m)
}
func (m *MechanismUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_MechanismUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_MechanismUpdate proto.InternalMessageInfo

func (m *MechanismUpdate) GetRemoteMechanisms() []*connection.Mechanism {
	if m != nil {
		return m.RemoteMechanisms
	}
	return nil
}

func (m *MechanismUpdate) GetLocalMechanisms() []*connection1.Mechanism {
	if m != nil {
		return m.LocalMechanisms
	}
	return nil
}

func init() {
	proto.RegisterType((*MechanismUpdate)(nil), "forwarder.MechanismUpdate")
}

func init() { proto.RegisterFile("forwarder.proto", fileDescriptor_19bff53f4d11db23) }

var fileDescriptor_19bff53f4d11db23 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x41, 0x4f, 0x32, 0x31,
	0x10, 0xcd, 0x7e, 0x5f, 0xd4, 0x50, 0x0f, 0x40, 0x0f, 0x86, 0x54, 0x0f, 0xc6, 0x93, 0xa7, 0xae,
	0xc1, 0xa3, 0x27, 0x43, 0x24, 0xf1, 0xc0, 0x85, 0xc4, 0xab, 0xa6, 0x94, 0x01, 0x1a, 0xdb, 0x4e,
	0x6d, 0x8b, 0x84, 0xdf, 0xe0, 0xdf, 0xf0, 0x87, 0x9a, 0xdd, 0x02, 0xbb, 0x18, 0xdc, 0x0b, 0x97,
	0xcd, 0xcc, 0x7b, 0xdb, 0xf7, 0x66, 0x5e, 0x4b, 0xda, 0x33, 0xf4, 0x2b, 0xe1, 0xa7, 0xe0, 0xb9,
	0xf3, 0x18, 0x91, 0xb6, 0x76, 0x00, 0x93, 0x73, 0x15, 0x17, 0xcb, 0x09, 0x97, 0x68, 0x72, 0x0b,
	0x71, 0x85, 0xfe, 0x3d, 0x80, 0xff, 0x54, 0x12, 0x0c, 0x84, 0xc5, 0x21, 0x48, 0xa2, 0x8d, 0x1e,
	0xb5, 0xd3, 0xc2, 0x42, 0x2e, 0x9c, 0xca, 0x35, 0x4a, 0xa1, 0x0b, 0xd8, 0x82, 0x8c, 0x0a, 0x6d,
	0xad, 0x4c, 0x7e, 0x6c, 0x7a, 0xbc, 0x89, 0x07, 0x83, 0x11, 0x1a, 0x5d, 0xc4, 0xf1, 0x2e, 0xd2,
	0x63, 0x08, 0x1b, 0xe1, 0xbd, 0x66, 0x63, 0xd1, 0x73, 0x71, 0xed, 0x20, 0xe4, 0x60, 0x5c, 0x5c,
	0xa7, 0x6f, 0x62, 0x6e, 0xbe, 0x33, 0xd2, 0x1e, 0x81, 0x5c, 0x08, 0xab, 0x82, 0x79, 0x71, 0x53,
	0x11, 0x81, 0x3e, 0x93, 0x6e, 0x1a, 0xfb, 0xcd, 0x6c, 0x99, 0xd0, 0xcb, 0xae, 0xff, 0xdf, 0x9e,
	0xf7, 0xaf, 0x78, 0x62, 0x78, 0x6d, 0x8b, 0xdd, 0xf1, 0x71, 0x27, 0x91, 0x3b, 0x20, 0xd0, 0x21,
	0xe9, 0x94, 0x31, 0xd7, 0x95, 0xfe, 0x95, 0x4a, 0x97, 0xbc, 0x24, 0x0e, 0x0b, 0xb5, 0x4b, 0xae,
	0xd2, 0xe9, 0x7f, 0x65, 0xa4, 0x35, 0xdc, 0x5e, 0x3e, 0x7d, 0x24, 0x67, 0x63, 0xf8, 0x58, 0x42,
	0x88, 0x94, 0xf1, 0xbd, 0x75, 0x07, 0x45, 0x33, 0x48, 0x0d, 0x6b, 0xe0, 0xe8, 0x03, 0x39, 0x19,
	0x68, 0x0c, 0xd0, 0x28, 0x70, 0xc1, 0xe7, 0x88, 0x73, 0x0d, 0x29, 0xab, 0xc9, 0x72, 0xc6, 0x9f,
	0x8a, 0xe8, 0xfa, 0xaf, 0xa4, 0x5b, 0xcd, 0x36, 0x42, 0xab, 0x22, 0xfa, 0x22, 0xb5, 0x4d, 0x59,
	0xdb, 0xff, 0x0f, 0x05, 0xc6, 0x78, 0xf5, 0xb6, 0x7f, 0xc5, 0x7f, 0x97, 0x4d, 0x4e, 0xcb, 0xbf,
	0xef, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x04, 0x3d, 0xe9, 0xfc, 0x01, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ForwarderClient is the client API for Forwarder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ForwarderClient interface {
	Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error)
	Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error)
}

type forwarderClient struct {
	cc *grpc.ClientConn
}

func NewForwarderClient(cc *grpc.ClientConn) ForwarderClient {
	return &forwarderClient{cc}
}

func (c *forwarderClient) Request(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*crossconnect.CrossConnect, error) {
	out := new(crossconnect.CrossConnect)
	err := c.cc.Invoke(ctx, "/forwarder.Forwarder/Request", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forwarderClient) Close(ctx context.Context, in *crossconnect.CrossConnect, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/forwarder.Forwarder/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForwarderServer is the server API for Forwarder service.
type ForwarderServer interface {
	Request(context.Context, *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error)
	Close(context.Context, *crossconnect.CrossConnect) (*empty.Empty, error)
}

// UnimplementedForwarderServer can be embedded to have forward compatible implementations.
type UnimplementedForwarderServer struct {
}

func (*UnimplementedForwarderServer) Request(ctx context.Context, req *crossconnect.CrossConnect) (*crossconnect.CrossConnect, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Request not implemented")
}
func (*UnimplementedForwarderServer) Close(ctx context.Context, req *crossconnect.CrossConnect) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}

func RegisterForwarderServer(s *grpc.Server, srv ForwarderServer) {
	s.RegisterService(&_Forwarder_serviceDesc, srv)
}

func _Forwarder_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(crossconnect.CrossConnect)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwarderServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.Forwarder/Request",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwarderServer).Request(ctx, req.(*crossconnect.CrossConnect))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forwarder_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(crossconnect.CrossConnect)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForwarderServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forwarder.Forwarder/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForwarderServer).Close(ctx, req.(*crossconnect.CrossConnect))
	}
	return interceptor(ctx, in, info, handler)
}

var _Forwarder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.Forwarder",
	HandlerType: (*ForwarderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _Forwarder_Request_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Forwarder_Close_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "forwarder.proto",
}

// MechanismsMonitorClient is the client API for MechanismsMonitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MechanismsMonitorClient interface {
	MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (MechanismsMonitor_MonitorMechanismsClient, error)
}

type mechanismsMonitorClient struct {
	cc *grpc.ClientConn
}

func NewMechanismsMonitorClient(cc *grpc.ClientConn) MechanismsMonitorClient {
	return &mechanismsMonitorClient{cc}
}

func (c *mechanismsMonitorClient) MonitorMechanisms(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (MechanismsMonitor_MonitorMechanismsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MechanismsMonitor_serviceDesc.Streams[0], "/forwarder.MechanismsMonitor/MonitorMechanisms", opts...)
	if err != nil {
		return nil, err
	}
	x := &mechanismsMonitorMonitorMechanismsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MechanismsMonitor_MonitorMechanismsClient interface {
	Recv() (*MechanismUpdate, error)
	grpc.ClientStream
}

type mechanismsMonitorMonitorMechanismsClient struct {
	grpc.ClientStream
}

func (x *mechanismsMonitorMonitorMechanismsClient) Recv() (*MechanismUpdate, error) {
	m := new(MechanismUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MechanismsMonitorServer is the server API for MechanismsMonitor service.
type MechanismsMonitorServer interface {
	MonitorMechanisms(*empty.Empty, MechanismsMonitor_MonitorMechanismsServer) error
}

// UnimplementedMechanismsMonitorServer can be embedded to have forward compatible implementations.
type UnimplementedMechanismsMonitorServer struct {
}

func (*UnimplementedMechanismsMonitorServer) MonitorMechanisms(req *empty.Empty, srv MechanismsMonitor_MonitorMechanismsServer) error {
	return status.Errorf(codes.Unimplemented, "method MonitorMechanisms not implemented")
}

func RegisterMechanismsMonitorServer(s *grpc.Server, srv MechanismsMonitorServer) {
	s.RegisterService(&_MechanismsMonitor_serviceDesc, srv)
}

func _MechanismsMonitor_MonitorMechanisms_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MechanismsMonitorServer).MonitorMechanisms(m, &mechanismsMonitorMonitorMechanismsServer{stream})
}

type MechanismsMonitor_MonitorMechanismsServer interface {
	Send(*MechanismUpdate) error
	grpc.ServerStream
}

type mechanismsMonitorMonitorMechanismsServer struct {
	grpc.ServerStream
}

func (x *mechanismsMonitorMonitorMechanismsServer) Send(m *MechanismUpdate) error {
	return x.ServerStream.SendMsg(m)
}

var _MechanismsMonitor_serviceDesc = grpc.ServiceDesc{
	ServiceName: "forwarder.MechanismsMonitor",
	HandlerType: (*MechanismsMonitorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MonitorMechanisms",
			Handler:       _MechanismsMonitor_MonitorMechanisms_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "forwarder.proto",
}
