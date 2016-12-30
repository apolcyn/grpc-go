package main

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	testpb "google.golang.org/grpc/benchmark/grpc_testing"
	"google.golang.org/grpc/grpclog"
)

func setPayload(p *testpb.Payload, t testpb.PayloadType, size int) {
	if size < 0 {
		grpclog.Fatalf("Requested a response with invalid length %d", size)
	}
	body := make([]byte, size)
	switch t {
	case testpb.PayloadType_COMPRESSABLE:
	case testpb.PayloadType_UNCOMPRESSABLE:
		grpclog.Fatalf("PayloadType UNCOMPRESSABLE is not supported")
	default:
		grpclog.Fatalf("Unsupported payload type: %d", t)
	}
	p.Type = t
	p.Body = body
	return
}

func newPayload(t testpb.PayloadType, size int) *testpb.Payload {
	p := new(testpb.Payload)
	setPayload(p, t, size)
	return p
}

func main() {
	grpclog.Println("start")
	p := grpc.NewProtoCodec()
	var wg sync.WaitGroup
	var total int64
	var done bool
	for i := 0; i < 6400; i++ {
		wg.Add(1)
		go func() {
			var count int64
			for !done {
				s := &testpb.SimpleResponse{
					Payload: newPayload(testpb.PayloadType_COMPRESSABLE, int(0)),
				}
				b, _ := p.Marshal(s)
				count += 1
				p.Unmarshal(b, &testpb.SimpleResponse{})
			}
			atomic.AddInt64(&total, count)
			wg.Done()
		}()
	}

	time.Sleep(time.Second * 10)
	done = true
	grpclog.Println("done")
	wg.Wait()

	grpclog.Println("total is %d", total)
	grpclog.Println("done")
}
