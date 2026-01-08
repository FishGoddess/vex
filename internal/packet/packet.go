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
	magic    uint32
	ptype    PacketType
	length   uint32
	sequence uint64
	data     []byte
}

// New returns a packet with given values.
func New(ptype PacketType, sequence uint64, data []byte) Packet {
	return Packet{
		magic:    magic,
		ptype:    ptype,
		length:   uint32(len(data)),
		sequence: sequence,
		data:     data,
	}
}
