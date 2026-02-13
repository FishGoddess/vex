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
	magic       = 1997811915
	headerBytes = 24
)

var (
	maxDataBytes = uint32(1<<32 - 1) // 4GB
)

var (
	errWrongMagic   = errors.New("vex: magic is wrong")
	errWrongLength  = errors.New("vex: length is wrong")
	errDataTooLarge = errors.New("vex: data is too large")
)

// ReadPacket reads a packet from reader and returns an error if failed.
func ReadPacket(reader io.Reader) (packet Packet, err error) {
	header := make([]byte, headerBytes)

	_, err = io.ReadFull(reader, header)
	if err != nil {
		return packet, err
	}

	endian := binary.BigEndian
	packet.id = endian.Uint64(header[0:8])
	packet.magic = endian.Uint32(header[8:12])
	packet.flags = endian.Uint64(header[12:20])
	packet.length = endian.Uint32(header[20:24])

	if packet.magic != magic {
		return packet, errWrongMagic
	}

	if packet.length <= 0 {
		return packet, nil
	}

	if packet.length > maxDataBytes {
		return packet, errDataTooLarge
	}

	packet.data = make([]byte, packet.length)

	_, err = io.ReadFull(reader, packet.data)
	return packet, err
}

// WritePacket writes a packet to writer and returns an error if failed.
func WritePacket(writer io.Writer, packet Packet) (err error) {
	if packet.magic != magic {
		return errWrongMagic
	}

	length := uint32(len(packet.data))
	if packet.length != length {
		return errWrongLength
	}

	if packet.length > maxDataBytes {
		return errDataTooLarge
	}

	endian := binary.BigEndian
	packetBytes := make([]byte, 0, headerBytes+packet.length)
	packetBytes = endian.AppendUint64(packetBytes, packet.id)
	packetBytes = endian.AppendUint32(packetBytes, packet.magic)
	packetBytes = endian.AppendUint64(packetBytes, packet.flags)
	packetBytes = endian.AppendUint32(packetBytes, packet.length)
	packetBytes = append(packetBytes, packet.data...)

	_, err = writer.Write(packetBytes)
	return err
}
