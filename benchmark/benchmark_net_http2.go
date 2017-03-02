package benchmark

import (
	"crypto/tls"
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"golang.org/x/net/http2"
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
	if uint32(len(body)) != length {
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

func grpcPingPongFunc(expectedReqSize int32, respSize int32) func(w http.ResponseWriter, req *http.Request) {
	pingPong := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/grpc")
		w.Header().Set("Content-Length", "-1")
		w.(http.Flusher).Flush()
		w.Header().Set("Trailer:grpc-status", "0")
		header := make([]byte, 5)
		body := make([]byte, expectedReqSize)
		out := make([]byte, expectedReqSize+5)
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
	return pingPong
}

// starts up the server without blocking
// returns an address string and a close func
func StartGrpcUsingGoNetHttp2(port string, config *tls.Config, reqSize int32, respSize int32) (string, func()) {
	addr := "0.0.0.0:" + port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	handler := http.HandlerFunc(grpcPingPongFunc(reqSize, respSize))
	srv := &http.Server{Addr: addr, TLSConfig: config, Handler: handler}
	http2.ConfigureServer(srv, nil)
	tlsListener := tls.NewListener(l, config)
	go log.Fatal(srv.Serve(tlsListener))
	return l.Addr().String(), func() { srv.Close() }
}
