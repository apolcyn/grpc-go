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
package transport // import "google.golang.org/grpc/transport"

import "github.com/golang/protobuf/proto"

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
	OnNewStream() Codec
}

type CodecPerTransportCreator interface {
	OnNewTransport() CodecPerStreamCreator
}

// ProtoCodec is a Codec implementation with protobuf. It is the default codec for gRPC.
type ProtoCodec struct{}

func (ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (ProtoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (ProtoCodec) String() string {
	return "proto"
}

type ProtoCodecPerStreamCreator struct {
	codec Codec
}

func (c ProtoCodecPerStreamCreator) OnNewStream() Codec {
	return c.codec
}

type ProtoCodecPerTransportCreator struct {
	codec Codec
}

func (c ProtoCodecPerTransportCreator) OnNewTransport() CodecPerStreamCreator {
	return &ProtoCodecPerStreamCreator{
		codec: c.codec,
	}
}

func NewProtoCodecPerTransportCreator(codec Codec) CodecPerTransportCreator {
	return &ProtoCodecPerTransportCreator{
		codec: codec,
	}
}

/***************/
type GenericCodecPerStreamCreator struct {
	codec Codec
}

func (c GenericCodecPerStreamCreator) OnNewStream() Codec {
	return c.codec
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
