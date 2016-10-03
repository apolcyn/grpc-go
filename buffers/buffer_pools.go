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

// Package metadata define the structure of the metadata supported by gRPC library.
// import "google.golang.org/grpc/buffers"
package buffers

import (
	"strconv"
	"sync"
)

var GlobalProtoBufferPool = NewProtobufBufferPool()

func getCreator(minCap uint) func() interface{} {
	return func() interface{} {
		return make([]byte, minCap)
	}
}

type BufferPool interface {
	GetBuf(size uint) []byte
	PutBuf([]byte)
}

var sizesCount = 20

type protobufBufferPool struct {
	bufferPoolSizes []uint
	bufferPools     []*sync.Pool
	poolMu          *sync.Mutex
}

func NewProtobufBufferPool() BufferPool {
	sizes = make([]int, sizesCount)
	for i := 0; i < len(sizes); i++ {
		sizes[i] = -1
	}

	return &protobufBufferPool{
		bufferPoolSizes: sizes,
		bufferPools:     make([]*sync.Pool, poolTableSize),
		poolMu:          new(sync.Mutex),
	}
}

func (p *protobufBufferPool) getPool(minCap uint) (int, *sync.Pool) {
	for i := 0; i < sizesCount; i++ {
		if p.bufferPoolSizes[i] == minCap {
			return i, p.bufferPools[i]
		}
		if p.bufferPoolSizes[i] == -1 {
			p.bufferPoolSizes[i] = minCap
			return i, nil
		}
	}
	panic("asking for too a pool size thats not cached" + strconv.FormatUint(uint64(minCap), 10))
	return -1, nil
}

func (p *protobufBufferPool) GetBuf(minCap uint) []byte {
	defer p.poolMu.Unlock()
	p.poolMu.Lock()
	var index, pool = p.getPool(minCap)
	if pool == nil {
		p.bufferPools[index] = &sync.Pool{New: getCreator(minCap)}
	}
	var out = p.bufferPools[index].Get().([]byte)
	if len(out) < int(minCap) {
		panic("something is wrong with buffer pool")
	}
	return out[0:minCap]
}

func (p *protobufBufferPool) PutBuf(buf []byte) {
	defer p.poolMu.Unlock()
	p.poolMu.Lock()
	var minCap = uint(len(buf))
	var index, pool = p.getPool(minCap)
	if pool == nil {
		p.bufferPools[index] = &sync.Pool{New: getCreator(minCap)}
	}
	p.bufferPools[index].Put(buf)
}

type noCacheBufferPool struct{}

func NewDefaultBufferPool() BufferPool {
	return &noCacheBufferPool{}
}

func (p *noCacheBufferPool) GetBuf(minCap uint) []byte {
	return make([]byte, minCap)
}

func (p *noCacheBufferPool) PutBuf([]byte) {
	return
}
