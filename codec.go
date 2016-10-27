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

type CodecPerStreamCreator interface {
	CreateCodec() Codec
	CollectCodec(Codec)
}

type CodecPerTransportCreator interface {
	OnNewTransport() CodecPerStreamCreator
}

// ProtoCodec is a Codec implementation with protobuf. It is the default codec for gRPC.
type ProtoCodec struct {
	unmarshalBuffer *proto.Buffer
	marshalBuffer   *marshalBuffer
}

func (p ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	var protoMsg = v.(proto.Message)
	var sizeNeeded = proto.Size(protoMsg)
	var bufToUse []byte
	mb := p.marshalBuffer
	buffer := mb.buffer

	if mb.lastBuf != nil && sizeNeeded <= len(mb.lastBuf) {
		bufToUse = mb.lastBuf
	} else {
		bufToUse = make([]byte, sizeNeeded)
	}
	buffer.SetBuf(bufToUse)
	buffer.Reset()
	err := buffer.Marshal(protoMsg)
	if err != nil {
		panic("something went wrong with unmarshaling")
	}
	out := buffer.Bytes()
	buffer.SetBuf(nil)
	mb.lastBuf = bufToUse
	return out, err
}

func (p ProtoCodec) Unmarshal(data []byte, v interface{}) error {
	buffer := p.unmarshalBuffer
	buffer.SetBuf(data)
	err := buffer.Unmarshal(v.(proto.Message))
	buffer.SetBuf(nil)
	return err
}

func (ProtoCodec) String() string {
	return "proto"
}

type ProtoCodecPerStreamCreator struct {
	protoCodecPool *sync.Pool
}

func (c ProtoCodecPerStreamCreator) CreateCodec() Codec {
	codec := c.protoCodecPool.Get().(Codec)
	if codec == nil {
		codec = &ProtoCodec{
			marshalBuffer:   newMarshalBuffer(),
			unmarshalBuffer: &proto.Buffer{},
		}
	}
	return codec
}

func (p ProtoCodecPerStreamCreator) CollectCodec(c Codec) {
	p.protoCodecPool.Put(c)
}

type ProtoCodecPerTransportCreator struct {
}

func (c ProtoCodecPerTransportCreator) OnNewTransport() CodecPerStreamCreator {
	protoCodecPool := &sync.Pool{
		New: func() interface{} {
			return &ProtoCodec{
				unmarshalBuffer: &proto.Buffer{},
				marshalBuffer:   newMarshalBuffer(),
			}
		},
	}

	return &ProtoCodecPerStreamCreator{
		protoCodecPool: protoCodecPool,
	}
}

func NewProtoCodecPerTransportCreator() CodecPerTransportCreator {
	return &ProtoCodecPerTransportCreator{}
}

/***************/
type GenericCodecPerStreamCreator struct {
	codec Codec
}

func (c GenericCodecPerStreamCreator) CreateCodec() Codec {
	return c.codec
}

func (c GenericCodecPerStreamCreator) CollectCodec(Codec) {
}

type GenericCodecPerTransportCreator struct {
	codec Codec
}

func (c GenericCodecPerTransportCreator) OnNewTransport() CodecPerStreamCreator {
	return &GenericCodecPerStreamCreator{
		codec: c.codec,
	}
}

func NewGenericCodecPerTransportCreator(codec Codec) CodecPerTransportCreator {
	return &GenericCodecPerTransportCreator{
		codec: codec,
	}
}

type marshalBuffer struct {
	buffer  *proto.Buffer
	lastBuf []byte
}

func newMarshalBuffer() *marshalBuffer {
	return &marshalBuffer{
		buffer: &proto.Buffer{},
	}
}
