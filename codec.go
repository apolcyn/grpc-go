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

/*
Package transport defines and implements message oriented communication channel
to complete various transactions (e.g., an RPC).
*/
package grpc

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/grpclog"
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

type codecCreator interface {
	// Provides a new stream with a codec to be used for it's lifetime
	CreateCodec() Codec
	// Returns a Codec to its pool
	CollectCodec(Codec)
}

type codecManagerCreator interface {
	// Provides a codecCreator to be used by a connection/transport.
	// This can control the scope of codec pools, e.g. global, per-conn, none
	onNewTransport() codecCreator
}

// protoCodec is a Codec implementation with protobuf. It is the default codec for gRPC.
type protoCodec struct {
	unmarshalBuffer *proto.Buffer
	marshalBuffer   *marshalBuffer
}

func (p protoCodec) Marshal(v interface{}) ([]byte, error) {
	var protoMsg = v.(proto.Message)
	var sizeNeeded = proto.Size(protoMsg)
	var currentSlice []byte

	mb := p.marshalBuffer
	buffer := mb.buffer

	if mb.lastSlice != nil && sizeNeeded <= len(mb.lastSlice) {
		currentSlice = mb.lastSlice
	} else {
		currentSlice = make([]byte, sizeNeeded)
	}
	buffer.SetBuf(currentSlice)
	buffer.Reset()
	err := buffer.Marshal(protoMsg)
	if err != nil {
		return nil, err
	}
	out := buffer.Bytes()
	buffer.SetBuf(nil)
	mb.lastSlice = currentSlice
	return out, err
}

func (p protoCodec) Unmarshal(data []byte, v interface{}) error {
	buffer := p.unmarshalBuffer
	buffer.SetBuf(data)
	err := buffer.Unmarshal(v.(proto.Message))
	buffer.SetBuf(nil)
	return err
}

func (protoCodec) String() string {
	return "proto"
}

// The protoCodec per-stream creators, and per-transport creator "managers"
// are meant to handle pooling of protoCodec structs and their buffers.
// The current goal is to keep a pool of buffers per transport connection,
// to be used on its streams.
type protoCodecCreator struct {
	protoCodecPool *sync.Pool
}

func (c protoCodecCreator) CreateCodec() Codec {
	codec := c.protoCodecPool.Get().(Codec)
	if codec == nil {
		codec = &protoCodec{
			marshalBuffer:   newMarshalBuffer(),
			unmarshalBuffer: &proto.Buffer{},
		}
	}
	return codec
}

func (p protoCodecCreator) CollectCodec(c Codec) {
	if c != nil {
		p.protoCodecPool.Put(c)
	}
	grpclog.Println("received nil codec in collect function")
}

type protoCodecManagerCreator struct {
}

// Called when a new connection is made. Sets up the pool to be used by
// that connection, for its streams
func (c protoCodecManagerCreator) onNewTransport() codecCreator {
	protoCodecPool := &sync.Pool{
		New: func() interface{} {
			return &protoCodec{
				unmarshalBuffer: &proto.Buffer{},
				marshalBuffer:   newMarshalBuffer(),
			}
		},
	}

	return &protoCodecCreator{
		protoCodecPool: protoCodecPool,
	}
}

func newProtoCodecManagerCreator() codecManagerCreator {
	return &protoCodecManagerCreator{}
}

// Keeps a buffer used for marshalling, and can also holds on to the last
// byte slice used for marshalling for reuse
type marshalBuffer struct {
	buffer    *proto.Buffer
	lastSlice []byte
}

func newMarshalBuffer() *marshalBuffer {
	return &marshalBuffer{
		buffer: &proto.Buffer{},
	}
}

// generic codec used with user-supplied codec.
// These are used in the same way as the default protoCodec managers,
// but result in the single user-supplied codec being used on every
// connection/stream.
type genericCodecCreator struct {
	codec Codec
}

func (c genericCodecCreator) CreateCodec() Codec {
	return c.codec
}

func (c genericCodecCreator) CollectCodec(Codec) {
}

type genericCodecManagerCreator struct {
	codec Codec
}

func (c genericCodecManagerCreator) onNewTransport() codecCreator {
	return &genericCodecCreator{
		codec: c.codec,
	}
}

func newGenericCodecManagerCreator(codec Codec) codecManagerCreator {
	return &genericCodecManagerCreator{
		codec: codec,
	}
}
