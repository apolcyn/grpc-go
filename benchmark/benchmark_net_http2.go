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

func ReadGrpcMsg(r io.Reader, header []byte, body []byte) error {
	if _, err := io.ReadFull(r, header); err != nil {
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			return io.EOF
		}
		log.Fatalf("read header: %v", err)
	}
	pf := header[0]
	if pf != 0 {
		log.Fatalf("expect not compressed")
	}
	length := binary.BigEndian.Uint32(header[1:])
	if len(body) != length {
		log.Fatalf("expected msg length %v; got %v", len(body), length)
	}
	if _, err := io.ReadFull(r, body); err != nil {
		log.Fatalf("read body: %v", err)
	}
	return nil
}

func WriteGrpcMsg(w io.Writer, body []byte, out []byte) {
	if len(out) != 5+len(body) {
		panic("invalid")
	}
	out[0] = 0
	binary.BigEndian.PutUint32(out[1:], uint32(len(body)))
	if n := copy(out[5:], body); n != len(body) {
		panic("bad copy")
	}
	if n, err := w.Write(out); err != nil || n != len(out) {
		log.Fatalf("write msg: %v", err)
	}
	w.(http.Flusher).Flush()
}

func grpcPingPongFunc(expectedMsgLen int) func(w http.ResponseWriter, req *http.Request) {
	pingPong := func(w http.ResponseWriter, req *http.Request) {
	        w.Header().Set("Content-Type", "application/grpc")
	        w.Header().Set("Content-Length", "-1")
	        w.(http.Flusher).Flush()
	        w.Header().Set("Trailer:grpc-status", "0")
	        header := make([]byte, 5)
	        body := make([]byte, expectedMsgLen)
	        out := make([]byte, expectedMsgLen+5)
	        for {
	        	if err := ReadGrpcMsg(req.Body, header, body); err != nil {
	        		if err == io.EOF {
	        			break
	        		}
	        		log.Fatal(err)
	        	}
	        	WriteGrpcMsg(w, body, out)
	        }
	        req.Body.Close()
	}
}

func StartGrpcUsingGoNetHttp2(port string, creds *tls.Config) {
	l, err := net.Listen("tcp", "0.0.0.0:" + port)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	handler := http.HandlerFunc(pingPong)
	if *useGrpc {
		log.Println("starting grpc")
		handler = http.HandlerFunc(grpcPingPong)
	}

	srv := &http.Server{Addr: "0.0.0.0:8080", TLSConfig: &config, Handler: http.HandlerFunc(grpcPingPongFunc(responseSize))}
	http2.ConfigureServer(srv, nil)
	tlsListener := tls.NewListener(l, &config)
	log.Fatal(srv.Serve(tlsListener))
}
