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
	Magic       = 0xC638B // 811915
	HeaderBytes = 16      // magic + type + length + sequence
)

var (
	ErrWrongMagic  = errors.New("vex: magic is wrong")
	ErrWrongLength = errors.New("vex: length is wrong")
)

// Decode decodes the packet from reader and returns an error if failed.
func Decode(reader io.Reader) (packet Packet, err error) {
	var header [HeaderBytes]byte

	_, err = io.ReadFull(reader, header[:])
	if err != nil {
		return packet, err
	}

	endian := binary.BigEndian
	magicAndType := endian.Uint32(header[0:4])
	packet.Magic = magicAndType >> 8
	packet.Type = PacketType(magicAndType & 0xFF)
	packet.Length = endian.Uint32(header[4:8])
	packet.Sequence = endian.Uint64(header[8:16])

	if packet.Magic != Magic {
		return packet, ErrWrongMagic
	}

	if packet.Length == 0 {
		return packet, nil
	}

	packet.Data = make([]byte, packet.Length)

	_, err = io.ReadFull(reader, packet.Data)
	return packet, err
}

// Encode encodes the packet to writer and returns an error if failed.
func Encode(writer io.Writer, packet Packet) (err error) {
	if packet.Magic != Magic {
		return ErrWrongMagic
	}

	length := uint32(len(packet.Data))
	if packet.Length != length {
		return ErrWrongLength
	}

	endian := binary.BigEndian
	magicAndType := (packet.Magic << 8) | uint32(packet.Type)
	packetBytes := make([]byte, 0, HeaderBytes+packet.Length)
	packetBytes = endian.AppendUint32(packetBytes, magicAndType)
	packetBytes = endian.AppendUint32(packetBytes, packet.Length)
	packetBytes = endian.AppendUint64(packetBytes, packet.Sequence)
	packetBytes = append(packetBytes, packet.Data...)

	_, err = writer.Write(packetBytes)
	return err
}
