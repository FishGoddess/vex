// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

const (
	PacketTypeRequest  PacketType = 1
	PacketTypeResponse PacketType = 2
	PacketTypeError    PacketType = 3
)

type PacketType = uint8

type Packet struct {
	Magic    uint32
	Type     PacketType
	Length   uint32
	Sequence uint64
	Data     []byte
}

// With sets the length and data to packet.
func (p *Packet) With(data []byte) {
	p.Length = uint32(len(data))
	p.Data = data
}
