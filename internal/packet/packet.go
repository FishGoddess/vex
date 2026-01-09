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

// Type returns the type of packet.
func (p *Packet) Type() PacketType {
	return p.ptype
}

// Data returns the data of packet.
func (p *Packet) Data() []byte {
	return p.data
}

// SetType sets the type of packet.
func (p *Packet) SetType(ptype PacketType) {
	p.ptype = ptype
}

// SetType sets the type of packet.
func (p *Packet) SetData(data []byte) {
	p.length = uint32(len(data))
	p.data = data
}
