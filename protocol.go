// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	magicSize   = 4                                                // 4 Byte
	versionSize = 1                                                // 1 Byte
	typeSize    = 1                                                // 1 Byte
	maxBodySize = 4                                                // 4 Byte
	headerSize  = magicSize + versionSize + typeSize + maxBodySize // 10 Byte

	magicNumber     = 0x755DD8C // Magic number is 123067788, for checking data package.
	protocolVersion = 1         // v1
)

const (
	packetTypeOK  PacketType = 0
	packetTypeErr PacketType = 1
)

var (
	endian = binary.BigEndian // All encodes/decodes between number and bytes use this endian.

	errMagicMismatch       = errors.New("vex: magic number in protocol doesn't match")
	errProtocolMismatch    = errors.New("vex: protocol between client and server doesn't match")
	errReadSizeMismatch    = errors.New("vex: read size less than expected size")
	errWrittenSizeMismatch = errors.New("vex: written size less than expected size")
)

// PacketType is the type of packet.
type PacketType = byte

func readPacketHeader(reader io.Reader) (packetType PacketType, bodySize int32, err error) {
	var header [headerSize]byte

	n, err := reader.Read(header[:])
	if err != nil {
		return 0, 0, err
	}

	if n != headerSize {
		return 0, 0, errReadSizeMismatch
	}

	magic := int32(endian.Uint32(header[:magicSize]))
	if magic != magicNumber {
		return 0, 0, errMagicMismatch
	}

	index := magicSize
	if header[index] != protocolVersion {
		return 0, 0, errProtocolMismatch
	}

	index += versionSize
	packetType = header[index]

	index += typeSize
	bodySize = int32(endian.Uint32(header[index:headerSize]))
	return packetType, bodySize, nil
}

func readPacketBody(reader io.Reader, bodySize int32) (body []byte, err error) {
	body = MakeBytes(bodySize) // May exceed if body size is too big.

	n, err := reader.Read(body)
	if err != nil {
		return nil, err
	}

	if n != int(bodySize) {
		return nil, errReadSizeMismatch
	}

	return body, nil
}

func readPacket(reader io.Reader) (packetType PacketType, body []byte, err error) {
	packetType, bodySize, err := readPacketHeader(reader)
	if err != nil {
		return 0, nil, err
	}

	if bodySize > 0 {
		body, err = readPacketBody(reader, bodySize)
	}

	return packetType, body, err
}

func writePacketHeader(writer io.Writer, packetType PacketType, bodySize int32) (err error) {
	var header [headerSize]byte
	endian.PutUint32(header[:magicSize], magicNumber)

	index := magicSize
	header[index] = protocolVersion

	index += versionSize
	header[index] = packetType

	index += typeSize
	binary.BigEndian.PutUint32(header[index:headerSize], uint32(bodySize))

	n, err := writer.Write(header[:])
	if err != nil {
		return err
	}

	if n != headerSize {
		return errWrittenSizeMismatch
	}

	return nil
}

func writePacketBody(writer io.Writer, body []byte) (err error) {
	n, err := writer.Write(body)
	if err != nil {
		return err
	}

	if n != len(body) {
		return errWrittenSizeMismatch
	}

	return nil
}

func writePacket(writer io.Writer, packetType PacketType, body []byte) (err error) {
	bodySize := int32(len(body))

	err = writePacketHeader(writer, packetType, bodySize)
	if err != nil {
		return err
	}

	if bodySize > 0 {
		err = writePacketBody(writer, body)
	}

	return nil
}
