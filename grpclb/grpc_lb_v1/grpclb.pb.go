// Code generated by protoc-gen-go.
// source: grpclb.proto
// DO NOT EDIT!

/*
Package grpc_lb_v1 is a generated protocol buffer package.

It is generated from these files:
	grpclb.proto

It has these top-level messages:
	Duration
	LoadBalanceRequest
	InitialLoadBalanceRequest
	ClientStats
	LoadBalanceResponse
	InitialLoadBalanceResponse
	ServerList
	Server
*/
package grpc_lb_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	grpc "google.golang.org/grpc"
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

type Duration struct {
	// Signed seconds of the span of time. Must be from -315,576,000,000
	// to +315,576,000,000 inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	// Signed fractions of a second at nanosecond resolution of the span
	// of time. Durations less than one second are represented with a 0
	// `seconds` field and a positive or negative `nanos` field. For durations
	// of one second or more, a non-zero value for the `nanos` field must be
	// of the same sign as the `seconds` field. Must be from -999,999,999
	// to +999,999,999 inclusive.
	Nanos int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Duration) Reset()                    { *m = Duration{} }
func (m *Duration) String() string            { return proto.CompactTextString(m) }
func (*Duration) ProtoMessage()               {}
func (*Duration) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Duration) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Duration) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

type LoadBalanceRequest struct {
	// Types that are valid to be assigned to LoadBalanceRequestType:
	//	*LoadBalanceRequest_InitialRequest
	//	*LoadBalanceRequest_ClientStats
	LoadBalanceRequestType isLoadBalanceRequest_LoadBalanceRequestType `protobuf_oneof:"load_balance_request_type"`
}

func (m *LoadBalanceRequest) Reset()                    { *m = LoadBalanceRequest{} }
func (m *LoadBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*LoadBalanceRequest) ProtoMessage()               {}
func (*LoadBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isLoadBalanceRequest_LoadBalanceRequestType interface {
	isLoadBalanceRequest_LoadBalanceRequestType()
}

type LoadBalanceRequest_InitialRequest struct {
	InitialRequest *InitialLoadBalanceRequest `protobuf:"bytes,1,opt,name=initial_request,json=initialRequest,oneof"`
}
type LoadBalanceRequest_ClientStats struct {
	ClientStats *ClientStats `protobuf:"bytes,2,opt,name=client_stats,json=clientStats,oneof"`
}

func (*LoadBalanceRequest_InitialRequest) isLoadBalanceRequest_LoadBalanceRequestType() {}
func (*LoadBalanceRequest_ClientStats) isLoadBalanceRequest_LoadBalanceRequestType()    {}

func (m *LoadBalanceRequest) GetLoadBalanceRequestType() isLoadBalanceRequest_LoadBalanceRequestType {
	if m != nil {
		return m.LoadBalanceRequestType
	}
	return nil
}

func (m *LoadBalanceRequest) GetInitialRequest() *InitialLoadBalanceRequest {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_InitialRequest); ok {
		return x.InitialRequest
	}
	return nil
}

func (m *LoadBalanceRequest) GetClientStats() *ClientStats {
	if x, ok := m.GetLoadBalanceRequestType().(*LoadBalanceRequest_ClientStats); ok {
		return x.ClientStats
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoadBalanceRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoadBalanceRequest_OneofMarshaler, _LoadBalanceRequest_OneofUnmarshaler, _LoadBalanceRequest_OneofSizer, []interface{}{
		(*LoadBalanceRequest_InitialRequest)(nil),
		(*LoadBalanceRequest_ClientStats)(nil),
	}
}

func _LoadBalanceRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoadBalanceRequest)
	// load_balance_request_type
	switch x := m.LoadBalanceRequestType.(type) {
	case *LoadBalanceRequest_InitialRequest:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitialRequest); err != nil {
			return err
		}
	case *LoadBalanceRequest_ClientStats:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ClientStats); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoadBalanceRequest.LoadBalanceRequestType has unexpected type %T", x)
	}
	return nil
}

func _LoadBalanceRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoadBalanceRequest)
	switch tag {
	case 1: // load_balance_request_type.initial_request
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(InitialLoadBalanceRequest)
		err := b.DecodeMessage(msg)
		m.LoadBalanceRequestType = &LoadBalanceRequest_InitialRequest{msg}
		return true, err
	case 2: // load_balance_request_type.client_stats
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ClientStats)
		err := b.DecodeMessage(msg)
		m.LoadBalanceRequestType = &LoadBalanceRequest_ClientStats{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoadBalanceRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoadBalanceRequest)
	// load_balance_request_type
	switch x := m.LoadBalanceRequestType.(type) {
	case *LoadBalanceRequest_InitialRequest:
		s := proto.Size(x.InitialRequest)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoadBalanceRequest_ClientStats:
		s := proto.Size(x.ClientStats)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type InitialLoadBalanceRequest struct {
	// Name of load balanced service (IE, service.grpc.gslb.google.com)
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *InitialLoadBalanceRequest) Reset()                    { *m = InitialLoadBalanceRequest{} }
func (m *InitialLoadBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*InitialLoadBalanceRequest) ProtoMessage()               {}
func (*InitialLoadBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *InitialLoadBalanceRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Contains client level statistics that are useful to load balancing. Each
// count should be reset to zero after reporting the stats.
type ClientStats struct {
	// The total number of requests sent by the client since the last report.
	TotalRequests int64 `protobuf:"varint,1,opt,name=total_requests,json=totalRequests" json:"total_requests,omitempty"`
	// The number of client rpc errors since the last report.
	ClientRpcErrors int64 `protobuf:"varint,2,opt,name=client_rpc_errors,json=clientRpcErrors" json:"client_rpc_errors,omitempty"`
	// The number of dropped requests since the last report.
	DroppedRequests int64 `protobuf:"varint,3,opt,name=dropped_requests,json=droppedRequests" json:"dropped_requests,omitempty"`
}

func (m *ClientStats) Reset()                    { *m = ClientStats{} }
func (m *ClientStats) String() string            { return proto.CompactTextString(m) }
func (*ClientStats) ProtoMessage()               {}
func (*ClientStats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ClientStats) GetTotalRequests() int64 {
	if m != nil {
		return m.TotalRequests
	}
	return 0
}

func (m *ClientStats) GetClientRpcErrors() int64 {
	if m != nil {
		return m.ClientRpcErrors
	}
	return 0
}

func (m *ClientStats) GetDroppedRequests() int64 {
	if m != nil {
		return m.DroppedRequests
	}
	return 0
}

type LoadBalanceResponse struct {
	// Types that are valid to be assigned to LoadBalanceResponseType:
	//	*LoadBalanceResponse_InitialResponse
	//	*LoadBalanceResponse_ServerList
	LoadBalanceResponseType isLoadBalanceResponse_LoadBalanceResponseType `protobuf_oneof:"load_balance_response_type"`
}

func (m *LoadBalanceResponse) Reset()                    { *m = LoadBalanceResponse{} }
func (m *LoadBalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*LoadBalanceResponse) ProtoMessage()               {}
func (*LoadBalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isLoadBalanceResponse_LoadBalanceResponseType interface {
	isLoadBalanceResponse_LoadBalanceResponseType()
}

type LoadBalanceResponse_InitialResponse struct {
	InitialResponse *InitialLoadBalanceResponse `protobuf:"bytes,1,opt,name=initial_response,json=initialResponse,oneof"`
}
type LoadBalanceResponse_ServerList struct {
	ServerList *ServerList `protobuf:"bytes,2,opt,name=server_list,json=serverList,oneof"`
}

func (*LoadBalanceResponse_InitialResponse) isLoadBalanceResponse_LoadBalanceResponseType() {}
func (*LoadBalanceResponse_ServerList) isLoadBalanceResponse_LoadBalanceResponseType()      {}

func (m *LoadBalanceResponse) GetLoadBalanceResponseType() isLoadBalanceResponse_LoadBalanceResponseType {
	if m != nil {
		return m.LoadBalanceResponseType
	}
	return nil
}

func (m *LoadBalanceResponse) GetInitialResponse() *InitialLoadBalanceResponse {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_InitialResponse); ok {
		return x.InitialResponse
	}
	return nil
}

func (m *LoadBalanceResponse) GetServerList() *ServerList {
	if x, ok := m.GetLoadBalanceResponseType().(*LoadBalanceResponse_ServerList); ok {
		return x.ServerList
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LoadBalanceResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LoadBalanceResponse_OneofMarshaler, _LoadBalanceResponse_OneofUnmarshaler, _LoadBalanceResponse_OneofSizer, []interface{}{
		(*LoadBalanceResponse_InitialResponse)(nil),
		(*LoadBalanceResponse_ServerList)(nil),
	}
}

func _LoadBalanceResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LoadBalanceResponse)
	// load_balance_response_type
	switch x := m.LoadBalanceResponseType.(type) {
	case *LoadBalanceResponse_InitialResponse:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.InitialResponse); err != nil {
			return err
		}
	case *LoadBalanceResponse_ServerList:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ServerList); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LoadBalanceResponse.LoadBalanceResponseType has unexpected type %T", x)
	}
	return nil
}

func _LoadBalanceResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LoadBalanceResponse)
	switch tag {
	case 1: // load_balance_response_type.initial_response
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(InitialLoadBalanceResponse)
		err := b.DecodeMessage(msg)
		m.LoadBalanceResponseType = &LoadBalanceResponse_InitialResponse{msg}
		return true, err
	case 2: // load_balance_response_type.server_list
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ServerList)
		err := b.DecodeMessage(msg)
		m.LoadBalanceResponseType = &LoadBalanceResponse_ServerList{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LoadBalanceResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LoadBalanceResponse)
	// load_balance_response_type
	switch x := m.LoadBalanceResponseType.(type) {
	case *LoadBalanceResponse_InitialResponse:
		s := proto.Size(x.InitialResponse)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LoadBalanceResponse_ServerList:
		s := proto.Size(x.ServerList)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type InitialLoadBalanceResponse struct {
	// This is an application layer redirect that indicates the client should use
	// the specified server for load balancing. When this field is non-empty in
	// the response, the client should open a separate connection to the
	// load_balancer_delegate and call the BalanceLoad method.
	LoadBalancerDelegate string `protobuf:"bytes,1,opt,name=load_balancer_delegate,json=loadBalancerDelegate" json:"load_balancer_delegate,omitempty"`
	// This interval defines how often the client should send the client stats
	// to the load balancer. Stats should only be reported when the duration is
	// positive.
	ClientStatsReportInterval *Duration `protobuf:"bytes,3,opt,name=client_stats_report_interval,json=clientStatsReportInterval" json:"client_stats_report_interval,omitempty"`
}

func (m *InitialLoadBalanceResponse) Reset()                    { *m = InitialLoadBalanceResponse{} }
func (m *InitialLoadBalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*InitialLoadBalanceResponse) ProtoMessage()               {}
func (*InitialLoadBalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *InitialLoadBalanceResponse) GetLoadBalancerDelegate() string {
	if m != nil {
		return m.LoadBalancerDelegate
	}
	return ""
}

func (m *InitialLoadBalanceResponse) GetClientStatsReportInterval() *Duration {
	if m != nil {
		return m.ClientStatsReportInterval
	}
	return nil
}

type ServerList struct {
	// Contains a list of servers selected by the load balancer. The list will
	// be updated when server resolutions change or as needed to balance load
	// across more servers. The client should consume the server list in order
	// unless instructed otherwise via the client_config.
	Servers []*Server `protobuf:"bytes,1,rep,name=servers" json:"servers,omitempty"`
	// Indicates the amount of time that the client should consider this server
	// list as valid. It may be considered stale after waiting this interval of
	// time after receiving the list. If the interval is not positive, the
	// client can assume the list is valid until the next list is received.
	ExpirationInterval *Duration `protobuf:"bytes,3,opt,name=expiration_interval,json=expirationInterval" json:"expiration_interval,omitempty"`
}

func (m *ServerList) Reset()                    { *m = ServerList{} }
func (m *ServerList) String() string            { return proto.CompactTextString(m) }
func (*ServerList) ProtoMessage()               {}
func (*ServerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ServerList) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

func (m *ServerList) GetExpirationInterval() *Duration {
	if m != nil {
		return m.ExpirationInterval
	}
	return nil
}

type Server struct {
	// A resolved address for the server, serialized in network-byte-order. It may
	// either be an IPv4 or IPv6 address.
	IpAddress []byte `protobuf:"bytes,1,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	// A resolved port number for the server.
	Port int32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	// An opaque but printable token given to the frontend for each pick. All
	// frontend requests for that pick must include the token in its initial
	// metadata. The token is used by the backend to verify the request and to
	// allow the backend to report load to the gRPC LB system.
	LoadBalanceToken string `protobuf:"bytes,3,opt,name=load_balance_token,json=loadBalanceToken" json:"load_balance_token,omitempty"`
	// Indicates whether this particular request should be dropped by the client
	// when this server is chosen from the list.
	DropRequest bool `protobuf:"varint,4,opt,name=drop_request,json=dropRequest" json:"drop_request,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Server) GetIpAddress() []byte {
	if m != nil {
		return m.IpAddress
	}
	return nil
}

func (m *Server) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Server) GetLoadBalanceToken() string {
	if m != nil {
		return m.LoadBalanceToken
	}
	return ""
}

func (m *Server) GetDropRequest() bool {
	if m != nil {
		return m.DropRequest
	}
	return false
}

func init() {
	proto.RegisterType((*Duration)(nil), "grpc.lb.v1.Duration")
	proto.RegisterType((*LoadBalanceRequest)(nil), "grpc.lb.v1.LoadBalanceRequest")
	proto.RegisterType((*InitialLoadBalanceRequest)(nil), "grpc.lb.v1.InitialLoadBalanceRequest")
	proto.RegisterType((*ClientStats)(nil), "grpc.lb.v1.ClientStats")
	proto.RegisterType((*LoadBalanceResponse)(nil), "grpc.lb.v1.LoadBalanceResponse")
	proto.RegisterType((*InitialLoadBalanceResponse)(nil), "grpc.lb.v1.InitialLoadBalanceResponse")
	proto.RegisterType((*ServerList)(nil), "grpc.lb.v1.ServerList")
	proto.RegisterType((*Server)(nil), "grpc.lb.v1.Server")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for LoadBalancer service

type LoadBalancerClient interface {
	// Bidirectional rpc to get a list of servers.
	BalanceLoad(ctx context.Context, opts ...grpc.CallOption) (LoadBalancer_BalanceLoadClient, error)
}

type loadBalancerClient struct {
	cc *grpc.ClientConn
}

func NewLoadBalancerClient(cc *grpc.ClientConn) LoadBalancerClient {
	return &loadBalancerClient{cc}
}

func (c *loadBalancerClient) BalanceLoad(ctx context.Context, opts ...grpc.CallOption) (LoadBalancer_BalanceLoadClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_LoadBalancer_serviceDesc.Streams[0], c.cc, "/grpc.lb.v1.LoadBalancer/BalanceLoad", opts...)
	if err != nil {
		return nil, err
	}
	x := &loadBalancerBalanceLoadClient{stream}
	return x, nil
}

type LoadBalancer_BalanceLoadClient interface {
	Send(*LoadBalanceRequest) error
	Recv() (*LoadBalanceResponse, error)
	grpc.ClientStream
}

type loadBalancerBalanceLoadClient struct {
	grpc.ClientStream
}

func (x *loadBalancerBalanceLoadClient) Send(m *LoadBalanceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *loadBalancerBalanceLoadClient) Recv() (*LoadBalanceResponse, error) {
	m := new(LoadBalanceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for LoadBalancer service

type LoadBalancerServer interface {
	// Bidirectional rpc to get a list of servers.
	BalanceLoad(LoadBalancer_BalanceLoadServer) error
}

func RegisterLoadBalancerServer(s *grpc.Server, srv LoadBalancerServer) {
	s.RegisterService(&_LoadBalancer_serviceDesc, srv)
}

func _LoadBalancer_BalanceLoad_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LoadBalancerServer).BalanceLoad(&loadBalancerBalanceLoadServer{stream})
}

type LoadBalancer_BalanceLoadServer interface {
	Send(*LoadBalanceResponse) error
	Recv() (*LoadBalanceRequest, error)
	grpc.ServerStream
}

type loadBalancerBalanceLoadServer struct {
	grpc.ServerStream
}

func (x *loadBalancerBalanceLoadServer) Send(m *LoadBalanceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *loadBalancerBalanceLoadServer) Recv() (*LoadBalanceRequest, error) {
	m := new(LoadBalanceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _LoadBalancer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.lb.v1.LoadBalancer",
	HandlerType: (*LoadBalancerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BalanceLoad",
			Handler:       _LoadBalancer_BalanceLoad_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpclb.proto",
}

func init() { proto.RegisterFile("grpclb.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 568 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x41, 0x6f, 0x13, 0x3d,
	0x10, 0xad, 0xbf, 0xb4, 0xfd, 0x9a, 0xd9, 0xd0, 0x84, 0x69, 0x54, 0x92, 0x50, 0x20, 0xac, 0x54,
	0x14, 0x50, 0x15, 0x20, 0x70, 0x01, 0x71, 0x21, 0xb4, 0x52, 0x2a, 0xf5, 0x80, 0x1c, 0x38, 0xaf,
	0x9c, 0x5d, 0xab, 0x5a, 0xb1, 0xac, 0x8d, 0xed, 0x46, 0x70, 0xe4, 0x0c, 0x7f, 0x84, 0x9f, 0x81,
	0xf8, 0x63, 0x68, 0x6d, 0x6f, 0x76, 0x4b, 0x14, 0xc1, 0x6d, 0xe7, 0x8d, 0xe7, 0xf9, 0x8d, 0xdf,
	0xcc, 0x42, 0xeb, 0x52, 0xc9, 0x38, 0x5b, 0x8c, 0xa5, 0x12, 0x46, 0x20, 0x14, 0xd1, 0x38, 0x5b,
	0x8c, 0x97, 0x4f, 0xc3, 0x97, 0xb0, 0x77, 0x7a, 0xa5, 0x98, 0x49, 0x45, 0x8e, 0x3d, 0xf8, 0x5f,
	0xf3, 0x58, 0xe4, 0x89, 0xee, 0x91, 0x21, 0x19, 0x35, 0x68, 0x19, 0x62, 0x17, 0x76, 0x72, 0x96,
	0x0b, 0xdd, 0xfb, 0x6f, 0x48, 0x46, 0x3b, 0xd4, 0x05, 0xe1, 0x4f, 0x02, 0x78, 0x21, 0x58, 0x32,
	0x65, 0x19, 0xcb, 0x63, 0x4e, 0xf9, 0xa7, 0x2b, 0xae, 0x0d, 0xbe, 0x85, 0x76, 0x9a, 0xa7, 0x26,
	0x65, 0x59, 0xa4, 0x1c, 0x64, 0xe9, 0x82, 0xc9, 0xf1, 0xb8, 0xba, 0x78, 0x7c, 0xee, 0x8e, 0xac,
	0xd7, 0xcf, 0xb6, 0xe8, 0xbe, 0xaf, 0x2f, 0x19, 0x5f, 0x41, 0x2b, 0xce, 0x52, 0x9e, 0x9b, 0x48,
	0x1b, 0x66, 0x9c, 0x8a, 0x60, 0x72, 0xab, 0x4e, 0xf7, 0xc6, 0xe6, 0xe7, 0x45, 0x7a, 0xb6, 0x45,
	0x83, 0xb8, 0x0a, 0xa7, 0xb7, 0xa1, 0x9f, 0x09, 0x96, 0x44, 0x0b, 0x77, 0x4d, 0x29, 0x2a, 0x32,
	0x5f, 0x24, 0x0f, 0x1f, 0x43, 0x7f, 0xa3, 0x12, 0x44, 0xd8, 0xce, 0xd9, 0x47, 0x6e, 0xe5, 0x37,
	0xa9, 0xfd, 0x0e, 0xbf, 0x11, 0x08, 0x6a, 0x97, 0xe1, 0x31, 0xec, 0x1b, 0x61, 0xaa, 0x5e, 0xcb,
	0xb7, 0xbb, 0x61, 0x51, 0xcf, 0xa4, 0xf1, 0x11, 0xdc, 0xf4, 0x2d, 0x28, 0x19, 0x47, 0x5c, 0x29,
	0xa1, 0x5c, 0x1f, 0x0d, 0xda, 0x76, 0x09, 0x2a, 0xe3, 0x33, 0x0b, 0xe3, 0x43, 0xe8, 0x24, 0x4a,
	0x48, 0xc9, 0x93, 0x8a, 0xb4, 0xe1, 0x8e, 0x7a, 0xbc, 0xa4, 0x0d, 0x7f, 0x11, 0x38, 0xb8, 0x26,
	0x5c, 0x4b, 0x91, 0x6b, 0x8e, 0x73, 0xe8, 0x54, 0x1e, 0x38, 0xcc, 0x9b, 0xf0, 0xe0, 0x6f, 0x26,
	0xb8, 0xd3, 0xb3, 0x2d, 0xda, 0x5e, 0xb9, 0xe0, 0x49, 0x5f, 0x40, 0xa0, 0xb9, 0x5a, 0x72, 0x15,
	0x65, 0xa9, 0x36, 0xde, 0x85, 0xc3, 0x3a, 0xdf, 0xdc, 0xa6, 0x2f, 0x52, 0xeb, 0x22, 0xe8, 0x55,
	0x34, 0x3d, 0x82, 0xc1, 0x1f, 0x1e, 0x38, 0x4e, 0x67, 0xc2, 0x0f, 0x02, 0x83, 0xcd, 0x52, 0xf0,
	0x39, 0x1c, 0xd6, 0x8b, 0x55, 0x94, 0xf0, 0x8c, 0x5f, 0x32, 0x53, 0x1a, 0xd3, 0xcd, 0xaa, 0x22,
	0x75, 0xea, 0x73, 0xf8, 0x1e, 0x8e, 0xea, 0x43, 0x13, 0x29, 0x2e, 0x85, 0x32, 0x51, 0x9a, 0x1b,
	0xae, 0x96, 0x2c, 0xb3, 0x2f, 0x1a, 0x4c, 0xba, 0x75, 0xf9, 0xe5, 0x26, 0xd0, 0x7e, 0x6d, 0x7e,
	0xa8, 0xad, 0x3b, 0xf7, 0x65, 0xe1, 0x57, 0x02, 0x50, 0xb5, 0x89, 0x27, 0xc5, 0xce, 0x14, 0x51,
	0xe1, 0x7b, 0x63, 0x14, 0x4c, 0x70, 0xfd, 0x3d, 0x68, 0x79, 0x04, 0xcf, 0xe0, 0x80, 0x7f, 0x96,
	0xa9, 0xbb, 0xe5, 0xdf, 0xa4, 0x60, 0x55, 0xb0, 0xd2, 0xf0, 0x9d, 0xc0, 0xae, 0xa3, 0xc6, 0x3b,
	0x00, 0xa9, 0x8c, 0x58, 0x92, 0x28, 0xae, 0xdd, 0xe8, 0xb5, 0x68, 0x33, 0x95, 0xaf, 0x1d, 0x50,
	0x4c, 0x70, 0xa1, 0xde, 0xef, 0xad, 0xfd, 0xc6, 0x13, 0xc0, 0x6b, 0x5e, 0x18, 0xf1, 0x81, 0xe7,
	0x56, 0x43, 0x93, 0x76, 0x6a, 0x4f, 0xf9, 0xae, 0xc0, 0xf1, 0x3e, 0xb4, 0x8a, 0xa1, 0x5b, 0xad,
	0xf2, 0xf6, 0x90, 0x8c, 0xf6, 0x68, 0x50, 0x60, 0x7e, 0x0a, 0x27, 0x0b, 0x68, 0xd5, 0x6c, 0x53,
	0x48, 0x21, 0xf0, 0xdf, 0x05, 0x8c, 0x77, 0xeb, 0x7d, 0xad, 0x6f, 0xd9, 0xe0, 0xde, 0xc6, 0xbc,
	0xf3, 0x7f, 0x44, 0x9e, 0x90, 0xc5, 0xae, 0xfd, 0x75, 0x3d, 0xfb, 0x1d, 0x00, 0x00, 0xff, 0xff,
	0x3e, 0x3b, 0xb5, 0x0e, 0xca, 0x04, 0x00, 0x00,
}
