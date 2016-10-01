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

package grpc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/transport"
	"google.golang.org/grpc/buffers"
)

// Codec defines the interface gRPC uses to encode and decode messages.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// String returns the name of the Codec implementation. The returned
	// string will be used as part of content type in transmission.
	String() string
}

// protoCodec is a Codec implementation with protobuf. It is the default codec for gRPC.
type protoCodec struct {
	readBuffer  *proto.Buffer
	writeBuffer *proto.Buffer
}

func NewProtoCodec() *protoCodec {
	return &protoCodec{
		readBuffer:  proto.NewBuffer(nil),
		writeBuffer: proto.NewBuffer(nil),
	}
}

func (c *protoCodec) Marshal(v interface{}) ([]byte, error) {
	var protoMsg = v.(proto.Message)
	var sizeNeeded = proto.Size(protoMsg)
	var buf = buffers.AllocBuffer(sizeNeeded)
	if len(buf) != sizeNeeded {
		panic("something is up here")
	}
	c.writeBuffer.SetBuf(buf)
	c.writeBuffer.Reset()
	err := c.writeBuffer.Marshal(protoMsg)
	if err != nil {
		panic("something went wrong with unmarshaling")
	}
	var out = c.writeBuffer.Bytes()
	c.writeBuffer.SetBuf(nil)
	return out, nil
}

func (c *protoCodec) Unmarshal(data []byte, v interface{}) error {
	c.readBuffer.SetBuf(data)
	err := c.readBuffer.Unmarshal(v.(proto.Message))
	c.readBuffer.SetBuf(nil)
	return err
}

func (protoCodec) String() string {
	return "proto"
}

// Compressor defines the interface gRPC uses to compress a message.
type Compressor interface {
	// Do compresses p into w.
	Do(w io.Writer, p []byte) error
	// Type returns the compression algorithm the Compressor uses.
	Type() string
}

// NewGZIPCompressor creates a Compressor based on GZIP.
func NewGZIPCompressor() Compressor {
	return &gzipCompressor{}
}

type gzipCompressor struct {
}

func (c *gzipCompressor) Do(w io.Writer, p []byte) error {
	z := gzip.NewWriter(w)
	if _, err := z.Write(p); err != nil {
		return err
	}
	return z.Close()
}

func (c *gzipCompressor) Type() string {
	return "gzip"
}

// Decompressor defines the interface gRPC uses to decompress a message.
type Decompressor interface {
	// Do reads the data from r and uncompress them.
	Do(r io.Reader) ([]byte, error)
	// Type returns the compression algorithm the Decompressor uses.
	Type() string
}

type gzipDecompressor struct {
}

// NewGZIPDecompressor creates a Decompressor based on GZIP.
func NewGZIPDecompressor() Decompressor {
	return &gzipDecompressor{}
}

func (d *gzipDecompressor) Do(r io.Reader) ([]byte, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	return ioutil.ReadAll(z)
}

func (d *gzipDecompressor) Type() string {
	return "gzip"
}

// callInfo contains all related configuration and information about an RPC.
type callInfo struct {
	failFast  bool
	headerMD  metadata.MD
	trailerMD metadata.MD
	traceInfo traceInfo // in trace.go
}

var defaultCallInfo = callInfo{failFast: true}

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {
	// before is called before the call is sent to any server.  If before
	// returns a non-nil error, the RPC fails with that error.
	before(*callInfo) error

	// after is called after the call has completed.  after cannot return an
	// error, so any failures should be reported via output parameters.
	after(*callInfo)
}

type beforeCall func(c *callInfo) error

func (o beforeCall) before(c *callInfo) error { return o(c) }
func (o beforeCall) after(c *callInfo)        {}

type afterCall func(c *callInfo)

func (o afterCall) before(c *callInfo) error { return nil }
func (o afterCall) after(c *callInfo)        { o(c) }

// Header returns a CallOptions that retrieves the header metadata
// for a unary RPC.
func Header(md *metadata.MD) CallOption {
	return afterCall(func(c *callInfo) {
		*md = c.headerMD
	})
}

// Trailer returns a CallOptions that retrieves the trailer metadata
// for a unary RPC.
func Trailer(md *metadata.MD) CallOption {
	return afterCall(func(c *callInfo) {
		*md = c.trailerMD
	})
}

// FailFast configures the action to take when an RPC is attempted on broken
// connections or unreachable servers. If failfast is true, the RPC will fail
// immediately. Otherwise, the RPC client will block the call until a
// connection is available (or the call is canceled or times out) and will retry
// the call if it fails due to a transient error. Please refer to
// https://github.com/grpc/grpc/blob/master/doc/fail_fast.md
func FailFast(failFast bool) CallOption {
	return beforeCall(func(c *callInfo) error {
		c.failFast = failFast
		return nil
	})
}

// The format of the payload: compressed or not?
type payloadFormat uint8

const (
	compressionNone payloadFormat = iota // no compression
	compressionMade
)

// recvMsg reads a complete gRPC message from the stream.
//
// It returns the message and its payload (compression/encoding)
// format. The caller owns the returned msg memory.
//
// If there is an error, possible values are:
//   * io.EOF, when no messages remain
//   * io.ErrUnexpectedEOF
//   * of type transport.ConnectionError
//   * of type transport.StreamError
// No other error values or types must be returned, which also means
// that the underlying io.Reader must not return an incompatible
// error.
func recvAndParseMsg(s *transport.Stream, maxMsgSize int) (pf payloadFormat, msg []byte, err error) {
	frame, err := s.ReadNextFrame()
	if err != nil {
		return 0, nil, err
	}

	// Parse the 5 byte header of a gRPC message. Find more detail
	// at http://www.grpc.io/docs/guides/wire.html.
	if len(frame) < 5 {
		return 0, nil, io.ErrUnexpectedEOF
	}
	pf = payloadFormat(frame[0])
	length := binary.BigEndian.Uint32(frame[1:5])

	if length == 0 {
		return pf, nil, nil
	}
	if length > uint32(maxMsgSize) {
		return 0, nil, Errorf(codes.Internal, "grpc: received message length %d exceeding the max size %d", length, maxMsgSize)
	}
	if length != uint32(len(frame[5:])) {
		return 0, nil, io.ErrUnexpectedEOF
	}

	return pf, frame[5:], nil
}

// encode serializes msg and prepends the message header. If msg is nil, it
// generates the message header of 0 message length.
func encode(c Codec, msg interface{}, cp Compressor, cbuf *bytes.Buffer) ([]byte, error) {
	var b []byte
	var length uint
	if msg != nil {
		var err error
		// TODO(zhaoq): optimize to reduce memory alloc and copying.
		b, err = c.Marshal(msg)
		if err != nil {
			return nil, err
		}
		if cp != nil {
			if err := cp.Do(cbuf, b); err != nil {
				return nil, err
			}
			b = cbuf.Bytes()
		}
		length = uint(len(b))
	}
	if length > math.MaxUint32 {
		return nil, Errorf(codes.InvalidArgument, "grpc: message too large (%d bytes)", length)
	}

	const (
		payloadLen = 1
		sizeLen    = 4
	)

	var buf = make([]byte, payloadLen+sizeLen+len(b))

	// Write payload format
	if cp == nil {
		buf[0] = byte(compressionNone)
	} else {
		buf[0] = byte(compressionMade)
	}
	// Write length of b into buf
	binary.BigEndian.PutUint32(buf[1:], uint32(length))
	// Copy encoded msg to buf
	copy(buf[5:], b)

	return buf, nil
}

func checkRecvPayload(pf payloadFormat, recvCompress string, dc Decompressor) error {
	switch pf {
	case compressionNone:
	case compressionMade:
		if dc == nil || recvCompress != dc.Type() {
			return Errorf(codes.Unimplemented, "grpc: Decompressor is not installed for grpc-encoding %q", recvCompress)
		}
	default:
		return Errorf(codes.Internal, "grpc: received unexpected payload format %d", pf)
	}
	return nil
}

func recv(c Codec, s *transport.Stream, dc Decompressor, m interface{}, maxMsgSize int) error {
	pf, d, err := recvAndParseMsg(s, maxMsgSize)
	if err != nil {
		return err
	}
	if err := checkRecvPayload(pf, s.RecvCompress(), dc); err != nil {
		return err
	}
	if pf == compressionMade {
		d, err = dc.Do(bytes.NewReader(d))
		if err != nil {
			return Errorf(codes.Internal, "grpc: failed to decompress the received message %v", err)
		}
	}
	if len(d) > maxMsgSize {
		// TODO: Revisit the error code. Currently keep it consistent with java
		// implementation.
		return Errorf(codes.Internal, "grpc: received a message of %d bytes exceeding %d limit", len(d), maxMsgSize)
	}
	if err := c.Unmarshal(d, m); err != nil {
		return Errorf(codes.Internal, "grpc: failed to unmarshal the received message %v", err)
	}
	buffers.FreeBuffer(d)
	return nil
}

// rpcError defines the status from an RPC.
type rpcError struct {
	code codes.Code
	desc string
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error: code = %d desc = %s", e.code, e.desc)
}

// Code returns the error code for err if it was produced by the rpc system.
// Otherwise, it returns codes.Unknown.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if e, ok := err.(*rpcError); ok {
		return e.code
	}
	return codes.Unknown
}

// ErrorDesc returns the error description of err if it was produced by the rpc system.
// Otherwise, it returns err.Error() or empty string when err is nil.
func ErrorDesc(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*rpcError); ok {
		return e.desc
	}
	return err.Error()
}

// Errorf returns an error containing an error code and a description;
// Errorf returns nil if c is OK.
func Errorf(c codes.Code, format string, a ...interface{}) error {
	if c == codes.OK {
		return nil
	}
	return &rpcError{
		code: c,
		desc: fmt.Sprintf(format, a...),
	}
}

// toRPCErr converts an error into a rpcError.
func toRPCErr(err error) error {
	switch e := err.(type) {
	case *rpcError:
		return err
	case transport.StreamError:
		return &rpcError{
			code: e.Code,
			desc: e.Desc,
		}
	case transport.ConnectionError:
		return &rpcError{
			code: codes.Internal,
			desc: e.Desc,
		}
	default:
		switch err {
		case context.DeadlineExceeded:
			return &rpcError{
				code: codes.DeadlineExceeded,
				desc: err.Error(),
			}
		case context.Canceled:
			return &rpcError{
				code: codes.Canceled,
				desc: err.Error(),
			}
		case ErrClientConnClosing:
			return &rpcError{
				code: codes.FailedPrecondition,
				desc: err.Error(),
			}
		}

	}
	return Errorf(codes.Unknown, "%v", err)
}

// convertCode converts a standard Go error into its canonical code. Note that
// this is only used to translate the error returned by the server applications.
func convertCode(err error) codes.Code {
	switch err {
	case nil:
		return codes.OK
	case io.EOF:
		return codes.OutOfRange
	case io.ErrClosedPipe, io.ErrNoProgress, io.ErrShortBuffer, io.ErrShortWrite, io.ErrUnexpectedEOF:
		return codes.FailedPrecondition
	case os.ErrInvalid:
		return codes.InvalidArgument
	case context.Canceled:
		return codes.Canceled
	case context.DeadlineExceeded:
		return codes.DeadlineExceeded
	}
	switch {
	case os.IsExist(err):
		return codes.AlreadyExists
	case os.IsNotExist(err):
		return codes.NotFound
	case os.IsPermission(err):
		return codes.PermissionDenied
	}
	return codes.Unknown
}

// SupportPackageIsVersion3 is referenced from generated protocol buffer files
// to assert that that code is compatible with this version of the grpc package.
//
// This constant may be renamed in the future if a change in the generated code
// requires a synchronised update of grpc-go and protoc-gen-go. This constant
// should not be referenced from any other code.
const SupportPackageIsVersion3 = true
