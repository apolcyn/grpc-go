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

import "google.golang.org/grpc/grpclog"

const (
	initialMsgBufSize = 1024
	// bufs[i] = cache of buffers with cap initialMsgBufSize<<i
	numBufSizes = 13
)

const profileAllocs = false

func slicealloc(n int) (buf []byte) {
	minCap := n
	for i := range slicePools.sliceCaches {
		if bufCap := initialMsgBufSize << uint(i); minCap <= bufCap {
			minCap = bufCap // round up to this buffer size class
			sliceCache := &slicePools.sliceCaches[i]
			buf = sliceCache.sliceAlloc()
			if buf != nil {
				if c := cap(buf); c < minCap {
					grpclog.Println("BUG: found buf size %d in buffer cache for size %d")
					break
				}
				buf = buf[0:n]
				return buf
			}
			break
		}
	}
	if minCap < initialMsgBufSize {
		minCap = initialMsgBufSize
	}
	buf = make([]byte, n, minCap)
	return buf
}

var bufSizeClass = make(map[int]int, numBufSizes) // (initialMsgBufSize << i) -> i

var slicePools struct {
	sliceCaches [numBufSizes]sliceCache
}

func init() {
	for i := 0; i < numBufSizes; i++ {
		bufSizeClass[initialMsgBufSize<<uint(i)] = i
	}
}

func slicefree(bufp *[]byte) {
	buf := *bufp
	if buf == nil {
		return
	}
	if i, ok := bufSizeClass[cap(buf)]; ok {
		slicePools.sliceCaches[i].sliceFree(buf)
	}
	*bufp = nil
}

type sliceCache struct {
	cache     *ringCache
	sliceSize int
}

func (sc *sliceCache) sliceAlloc() []byte {
	out := sc.cache.pop()
	if out == nil {
		out = make([]byte, sc.sliceSize)
	}
	return out.([]byte)
}

func (sc *sliceCache) sliceFree(slice []byte) {
	if slice == nil {
		panic("freeing a nil slice")
	}
	if cap(slice) != sc.sliceSize {
		panic("freeing a slice of incorrect length")
	}
	sc.cache.push(slice)
}
