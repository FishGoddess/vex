// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import "errors"

const (
	flagError = 0x1
)

type Packet struct {
	id     uint64
	magic  uint32
	flags  uint64
	length uint32
	data   []byte
}

// New returns a new packet with id.
func New(id uint64) Packet {
	packet := Packet{id: id, magic: magic}
	return packet
}

func (p *Packet) flagSet(flag uint64) bool {
	return (p.flags & flag) > 0
}

// ID returns the id of packet.
func (p *Packet) ID() uint64 {
	return p.id
}

// Data returns the data of packet and returns an error if it's an error packet.
func (p *Packet) Data() ([]byte, error) {
	if p.flagSet(flagError) {
		err := errors.New(string(p.data))
		return nil, err
	}

	return p.data, nil
}

func (p *Packet) setFlag(flag uint64) {
	p.flags = p.flags | flag
}

// SetData sets the data and its length to packet.
func (p *Packet) SetData(data []byte) {
	p.length = uint32(len(data))
	p.data = data
}

// SetError sets the error and an error flag to packet.
func (p *Packet) SetError(err error) {
	if err == nil {
		return
	}

	p.setFlag(flagError)
	p.SetData([]byte(err.Error()))
}
