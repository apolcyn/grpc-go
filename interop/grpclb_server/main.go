package main

import (
	"flag"
	"net"
	"strconv"
	"fmt"
	"encoding/binary"

	"google.golang.org/grpc"
	lbpb "google.golang.org/grpc/grpclb/grpc_lb_v1/service"
	"google.golang.org/grpc/grpclb/grpc_lb_v1/messages"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	backend_port = flag.Int("backend_port", 20000, "The backend server port")
)

type loadBalancerServer struct {
}

func (*loadBalancerServer) BalanceLoad(lb_server lbpb.LoadBalancer_BalanceLoadServer) error {
	print("Request msg recv begin")
	_, _ = lb_server.Recv()
	print("Request msg recv done")
	localhost_ipv4 := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(localhost_ipv4, 0x7f000001)
	server := &messages.Server{
		IpAddress: localhost_ipv4,
		Port: int32(*backend_port),
	}
	res := &messages.LoadBalanceResponse{
		LoadBalanceResponseType: &messages.LoadBalanceResponse_ServerList{
			ServerList: &messages.ServerList{
				Servers: []*messages.Server{server},
			},
		},
	}
	lb_server.Send(res)
	print("Request msg send done")
	return nil
}



func main() {
	flag.Parse()
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	fmt.Printf("begin listing on %d\n", *port)
	lis, _ := net.Listen("tcp", ":"+strconv.Itoa(*port))
	lbpb.RegisterLoadBalancerServer(server, &loadBalancerServer{})
	server.Serve(lis)
}
