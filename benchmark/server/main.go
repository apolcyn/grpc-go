package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"math"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/benchmark"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var (
	duration = flag.Int("duration", math.MaxInt32, "The duration in seconds to run the benchmark server")
)

type byteBufCodec struct {
}

func (byteBufCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.(*[]byte)
	if !ok {
		return nil, fmt.Errorf("failed to marshal: %v is not type of *[]byte", v)
	}
	return *b, nil
}

func (byteBufCodec) Unmarshal(data []byte, v interface{}) error {
	b, ok := v.(*[]byte)
	if !ok {
		return fmt.Errorf("failed to marshal: %v is not type of *[]byte", v)
	}
	*b = data
	return nil
}

func (byteBufCodec) String() string {
	return "bytebuffer"
}

const (
	byteBufRespSize = 2
)

func main() {
	flag.Parse()
	go func() {
		lis, err := net.Listen("tcp", ":0")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}
		grpclog.Println("Server profiling address: ", lis.Addr().String())
		if err := http.Serve(lis, nil); err != nil {
			grpclog.Fatalf("Failed to serve: %v", err)
		}
	}()
	config := &tls.Config{}
	certFile := "server.crt"
	keyFile := "server.key"
	config.NextProtos = []string{"h2"}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], _ = tls.LoadX509KeyPair(certFile, keyFile)
	creds := credentials.NewTLS(config)
	addr, stopper := benchmark.StartServer(benchmark.ServerInfo{Addr: ":8081", Type: "bytebuf", Metadata: int32(byteBufRespSize)}, grpc.CustomCodec(&byteBufCodec{}), grpc.Creds(creds)) // listen on all interfaces
	grpclog.Println("Server Address: ", addr)
	<-time.After(time.Duration(*duration) * time.Second)
	stopper()
}
