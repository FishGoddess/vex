// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import "sync"

const KB = 1024

var bufferPool = sync.Pool{
	New: func() any {
		return &buffer{
			data: make([]byte, 0, 1*KB),
		}
	},
}

type buffer struct {
	data []byte
}

func (b *buffer) Data(length int, capacity int) []byte {
	for i := len(b.data); i < length; i++ {
		b.data = append(b.data, 0)
	}

	for i := cap(b.data); i < capacity; i++ {
		b.data = append(b.data, 0)
	}

	b.data = b.data[:length]
	return b.data
}

func acquireBuffer() *buffer {
	buff := bufferPool.Get().(*buffer)
	return buff
}

func releaseBuffer(buff *buffer) bool {
	if cap(buff.data) > 16*KB {
		return false
	}

	bufferPool.Put(buff)
	return true
}
