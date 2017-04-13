/*
 *
 * Copyright 2014, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package test_client

import (
	"flag"
	"fmt"
	"math"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	anypb "github.com/golang/protobuf/ptypes/any"
	spb "github.com/google/go-genproto/googleapis/rpc/status"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
	testpb "google.golang.org/grpc/test/grpc_testing"
)

var (
	// For headers:
	testMetadata = metadata.MD{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	testMetadata2 = metadata.MD{
		"key1": []string{"value12"},
		"key2": []string{"value22"},
	}
	// For trailers:
	testTrailerMetadata = metadata.MD{
		"tkey1": []string{"trailerValue1"},
		"tkey2": []string{"trailerValue2"},
	}
	testTrailerMetadata2 = metadata.MD{
		"tkey1": []string{"trailerValue12"},
		"tkey2": []string{"trailerValue22"},
	}
	// capital "Key" is illegal in HTTP/2.
	malformedHTTP2Metadata = metadata.MD{
		"Key": []string{"foo"},
	}
	testAppUA     = "myApp1/1.0 myApp2/0.9"
	failAppUA     = "fail-this-RPC"
	detailedError = status.ErrorProto(&spb.Status{
		Code:    int32(codes.DataLoss),
		Message: "error for testing: " + failAppUA,
		Details: []*anypb.Any{{
			TypeUrl: "url",
			Value:   []byte{6, 0, 0, 6, 1, 3},
		}},
	})
)

var raceMode bool // set by race_test.go in race mode

type testServer struct {
	security           string // indicate the authentication protocol used by this server.
	earlyFail          bool   // whether to error out the execution of a service handler prematurely.
	setAndSendHeader   bool   // whether to call setHeader and sendHeader.
	setHeaderOnly      bool   // whether to only call setHeader, not sendHeader.
	multipleSetTrailer bool   // whether to call setTrailer multiple times.
}

func newPayload(t testpb.PayloadType, size int32) (*testpb.Payload, error) {
	if size < 0 {
		return nil, fmt.Errorf("Requested a response with invalid length %d", size)
	}
	body := make([]byte, size)
	switch t {
	case testpb.PayloadType_COMPRESSABLE:
	case testpb.PayloadType_UNCOMPRESSABLE:
		return nil, fmt.Errorf("PayloadType UNCOMPRESSABLE is not supported")
	default:
		return nil, fmt.Errorf("Unsupported payload type: %d", t)
	}
	return &testpb.Payload{
		Type: t.Enum(),
		Body: body,
	}, nil
}

func (s *testServer) EmptyCall(ctx context.Context, in *testpb.Empty) (*testpb.Empty, error) {
	return new(testpb.Empty), nil
}

func (s *testServer) UnaryCall(ctx context.Context, in *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
	md, ok := metadata.FromContext(ctx)
	if ok {
		if _, exists := md[":authority"]; !exists {
			return nil, grpc.Errorf(codes.DataLoss, "expected an :authority metadata: %v", md)
		}
		if s.setAndSendHeader {
			if err := grpc.SetHeader(ctx, md); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SetHeader(_, %v) = %v, want <nil>", md, err)
			}
			if err := grpc.SendHeader(ctx, testMetadata2); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SendHeader(_, %v) = %v, want <nil>", testMetadata2, err)
			}
		} else if s.setHeaderOnly {
			if err := grpc.SetHeader(ctx, md); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SetHeader(_, %v) = %v, want <nil>", md, err)
			}
			if err := grpc.SetHeader(ctx, testMetadata2); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SetHeader(_, %v) = %v, want <nil>", testMetadata2, err)
			}
		} else {
			if err := grpc.SendHeader(ctx, md); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SendHeader(_, %v) = %v, want <nil>", md, err)
			}
		}
		if err := grpc.SetTrailer(ctx, testTrailerMetadata); err != nil {
			return nil, grpc.Errorf(grpc.Code(err), "grpc.SetTrailer(_, %v) = %v, want <nil>", testTrailerMetadata, err)
		}
		if s.multipleSetTrailer {
			if err := grpc.SetTrailer(ctx, testTrailerMetadata2); err != nil {
				return nil, grpc.Errorf(grpc.Code(err), "grpc.SetTrailer(_, %v) = %v, want <nil>", testTrailerMetadata2, err)
			}
		}
	}
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.DataLoss, "failed to get peer from ctx")
	}
	if pr.Addr == net.Addr(nil) {
		return nil, grpc.Errorf(codes.DataLoss, "failed to get peer address")
	}
	if s.security != "" {
		// Check Auth info
		var authType, serverName string
		switch info := pr.AuthInfo.(type) {
		case credentials.TLSInfo:
			authType = info.AuthType()
			serverName = info.State.ServerName
		default:
			return nil, grpc.Errorf(codes.Unauthenticated, "Unknown AuthInfo type")
		}
		if authType != s.security {
			return nil, grpc.Errorf(codes.Unauthenticated, "Wrong auth type: got %q, want %q", authType, s.security)
		}
		if serverName != "x.test.youtube.com" {
			return nil, grpc.Errorf(codes.Unauthenticated, "Unknown server name %q", serverName)
		}
	}
	// Simulate some service delay.
	time.Sleep(time.Second)

	payload, err := newPayload(in.GetResponseType(), in.GetResponseSize())
	if err != nil {
		return nil, err
	}

	return &testpb.SimpleResponse{
		Payload: payload,
	}, nil
}

func (s *testServer) StreamingOutputCall(args *testpb.StreamingOutputCallRequest, stream testpb.TestService_StreamingOutputCallServer) error {
	return nil
}

func (s *testServer) StreamingInputCall(stream testpb.TestService_StreamingInputCallServer) error {
	return nil
}

func (s *testServer) FullDuplexCall(stream testpb.TestService_FullDuplexCallServer) error {
	return nil
}

func (s *testServer) HalfDuplexCall(stream testpb.TestService_HalfDuplexCallServer) error {
	return nil
}


const tlsDir = "testdata/"

type env struct {
	name        string
	network     string // The type of network such as tcp, unix, etc.
	security    string // The security protocol such as TLS, SSH, etc.
	httpHandler bool   // whether to use the http.Handler ServerTransport; requires TLS
	balancer    bool   // whether to use balancer
}

func (e env) runnable() bool {
	if runtime.GOOS == "windows" && e.network == "unix" {
		return false
	}
	return true
}

func (e env) dialer(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(e.network, addr, timeout)
}

var (
	tcpClearEnv   = env{name: "tcp-clear", network: "tcp", balancer: true}
	// TODO add handlerEnv back when ServeHTTP is stable.
	allEnv = []env{tcpClearEnv}
)

var onlyEnv = flag.String("only_env", "", "If non-empty, one of 'tcp-clear', 'tcp-tls', 'unix-clear', 'unix-tls', or 'handler-tls' to only run the tests for that environment. Empty means all.")

func listTestEnv() (envs []env) {
	if *onlyEnv != "" {
		for _, e := range allEnv {
			if e.name == *onlyEnv {
				if !e.runnable() {
					panic(fmt.Sprintf("--only_env environment %q does not run on %s", *onlyEnv, runtime.GOOS))
				}
				return []env{e}
			}
		}
		panic(fmt.Sprintf("invalid --only_env value %q", *onlyEnv))
	}
	for _, e := range allEnv {
		if e.runnable() {
			envs = append(envs, e)
		}
	}
	return envs
}

// test is an end-to-end test. It should be created with the newTest
// func, modified as needed, and then started with its startServer method.
// It should be cleaned up with the tearDown method.
type test struct {
	t *testing.T
	e env

	ctx    context.Context // valid for life of test, before tearDown
	cancel context.CancelFunc

	// Configurable knobs, after newTest returns:
	testServer        testpb.TestServiceServer // nil means none
	healthServer      *health.Server           // nil means disabled
	maxStream         uint32
	tapHandle         tap.ServerInHandle
	maxMsgSize        int
	userAgent         string
	clientCompression bool
	serverCompression bool
	unaryClientInt    grpc.UnaryClientInterceptor
	streamClientInt   grpc.StreamClientInterceptor
	unaryServerInt    grpc.UnaryServerInterceptor
	streamServerInt   grpc.StreamServerInterceptor
	unknownHandler    grpc.StreamHandler
	sc                <-chan grpc.ServiceConfig

	// srv and srvAddr are set once startServer is called.
	srv     *grpc.Server
	srvAddr string

	cc          *grpc.ClientConn // nil until requested via clientConn
	restoreLogs func()           // nil unless declareLogNoise is used
}

func (te *test) tearDown() {
	if te.cancel != nil {
		te.cancel()
		te.cancel = nil
	}
	if te.cc != nil {
		te.cc.Close()
		te.cc = nil
	}
	if te.restoreLogs != nil {
		te.restoreLogs()
		te.restoreLogs = nil
	}
	if te.srv != nil {
		te.srv.Stop()
	}
}

// newTest returns a new test using the provided testing.T and
// environment.  It is returned with default values. Tests should
// modify it before calling its startServer and clientConn methods.
func newTest(t *testing.T, e env) *test {
	te := &test{
		t:         t,
		e:         e,
		maxStream: math.MaxUint32,
	}
	te.ctx, te.cancel = context.WithCancel(context.Background())
	return te
}

// startServer starts a gRPC server listening. Callers should defer a
// call to te.tearDown to clean up.
func (te *test) startServer(ts testpb.TestServiceServer) {
	te.testServer = ts
	te.t.Logf("Running test in %s environment...", te.e.name)
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(te.maxStream)}
	if te.maxMsgSize > 0 {
		sopts = append(sopts, grpc.MaxMsgSize(te.maxMsgSize))
	}
	la := "localhost:0"
	lis, err := net.Listen(te.e.network, la)
	if err != nil {
		te.t.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(sopts...)
	te.srv = s
	testpb.RegisterTestServiceServer(s, te.testServer)
	_, port, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		te.t.Fatalf("Failed to parse listener address: %v", err)
	}
	addr := "localhost:" + port

	go s.Serve(lis)
	te.srvAddr = addr
}

func (te *test) clientConn() *grpc.ClientConn {
	if te.cc != nil {
		return te.cc
	}
	opts := []grpc.DialOption{
		grpc.WithDialer(te.e.dialer),
		grpc.WithUserAgent(te.userAgent),
	}

	if te.maxMsgSize > 0 {
		opts = append(opts, grpc.WithMaxMsgSize(te.maxMsgSize))
	}
	opts = append(opts, grpc.WithInsecure())
	var err error
	te.cc, err = grpc.Dial(te.srvAddr, opts...)
	if err != nil {
		te.t.Fatalf("Dial(%q) = %v", te.srvAddr, err)
	}
	return te.cc
}

func TestLargeUnary(t *testing.T) {
	for _, e := range listTestEnv() {
		testLargeUnary(t, e)
	}
}

func testLargeUnary(t *testing.T, e env) {
	te := newTest(t, e)
	te.startServer(&testServer{security: e.security})
	defer te.tearDown()
	tc := testpb.NewTestServiceClient(te.clientConn())

	const argSize = 271828
	const respSize = 314159

	payload, err := newPayload(testpb.PayloadType_COMPRESSABLE, argSize)
	if err != nil {
		t.Fatal(err)
	}

	req := &testpb.SimpleRequest{
		ResponseType: testpb.PayloadType_COMPRESSABLE.Enum(),
		ResponseSize: proto.Int32(respSize),
		Payload:      payload,
	}
	reply, err := tc.UnaryCall(context.Background(), req)
	if err != nil {
		t.Fatalf("TestService/UnaryCall(_, _) = _, %v, want _, <nil>", err)
	}
	pt := reply.GetPayload().GetType()
	ps := len(reply.GetPayload().GetBody())
	if pt != testpb.PayloadType_COMPRESSABLE || ps != respSize {
		t.Fatalf("Got the reply with type %d len %d; want %d, %d", pt, ps, testpb.PayloadType_COMPRESSABLE, respSize)
	}
}
