// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	magic       = 0xC638B // 811915
	headerBytes = 16      // magic + type + length + sequence
)

var (
	errWrongMagic  = errors.New("vex: magic is wrong")
	errWrongLength = errors.New("vex: length is wrong")
)

// Decode decodes the packet from reader and returns an error if failed.
func Decode(reader io.Reader) (packet Packet, err error) {
	var header [headerBytes]byte

	_, err = io.ReadFull(reader, header[:])
	if err != nil {
		return packet, err
	}

	endian := binary.BigEndian
	magicAndType := endian.Uint32(header[0:4])
	packet.magic = magicAndType >> 8
	packet.ptype = PacketType(magicAndType & 0xFF)
	packet.length = endian.Uint32(header[4:8])
	packet.sequence = endian.Uint64(header[8:16])

	if packet.magic != magic {
		return packet, errWrongMagic
	}

	if packet.length == 0 {
		return packet, nil
	}

	packet.data = make([]byte, packet.length)

	_, err = io.ReadFull(reader, packet.data)
	return packet, err
}

// Encode encodes the packet to writer and returns an error if failed.
func Encode(writer io.Writer, packet Packet) (err error) {
	if packet.magic != magic {
		return errWrongMagic
	}

	length := uint32(len(packet.data))
	if packet.length != length {
		return errWrongLength
	}

	endian := binary.BigEndian
	magicAndType := (packet.magic << 8) | uint32(packet.ptype)
	packetBytes := make([]byte, 0, headerBytes+packet.length)
	packetBytes = endian.AppendUint32(packetBytes, magicAndType)
	packetBytes = endian.AppendUint32(packetBytes, packet.length)
	packetBytes = endian.AppendUint64(packetBytes, packet.sequence)
	packetBytes = append(packetBytes, packet.data...)

	_, err = writer.Write(packetBytes)
	return err
}
