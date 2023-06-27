// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package packet

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	magicBits    = 24
	typeBits     = 8
	bodySizeBits = 32
	headerSize   = 8 // Bytes

	maxMagic    = 1<<magicBits - 1
	maxType     = 1<<typeBits - 1
	maxBodySize = 1<<bodySizeBits - 1

	// Ha! Guess what this number means?
	magicNumber = 0xC638B
)

const (
	typeOK  Type = 0
	typeErr Type = 1
)

var (
	// All encodes/decodes between number and bytes use this endian.
	endian = binary.BigEndian

	errMagicMismatch       = errors.New("vex: magic number in protocol doesn't match")
	errReadSizeMismatch    = errors.New("vex: read size less than expected size")
	errWrittenSizeMismatch = errors.New("vex: written size less than expected size")
)

// Type is the type of packet.
type Type = byte

func readPacketHeader(reader io.Reader) (Type, int32, error) {
	var headerBytes [headerSize]byte

	n, err := reader.Read(headerBytes[:])
	if err != nil {
		return 0, 0, err
	}

	if n != headerSize {
		return 0, 0, errReadSizeMismatch
	}

	header := endian.Uint64(headerBytes[:])
	magic := (header >> (typeBits + bodySizeBits)) & maxMagic

	if magic != magicNumber {
		return 0, 0, errMagicMismatch
	}

	packetType := Type((header >> bodySizeBits) & maxType)
	bodySize := int32(header & maxBodySize)
	return packetType, bodySize, nil
}

func readPacketBody(reader io.Reader, bodySize int32) ([]byte, error) {
	// Memory may exceed if body size is too big.
	body := make([]byte, bodySize, bodySize)

	n, err := io.ReadFull(reader, body)
	if err == io.ErrUnexpectedEOF {
		return nil, errReadSizeMismatch
	}

	if err != nil {
		return nil, err
	}

	if n != int(bodySize) {
		return nil, errReadSizeMismatch
	}

	return body, nil
}

func readPacket(reader io.Reader) (Type, []byte, error) {
	packetType, bodySize, err := readPacketHeader(reader)
	if err != nil {
		return 0, nil, err
	}

	var body []byte
	if bodySize > 0 {
		body, err = readPacketBody(reader, bodySize)
	}

	return packetType, body, err
}

func writePacketHeader(writer io.Writer, packetType Type, bodySize int32) error {
	var headerBytes [headerSize]byte
	var header = magicNumber<<(typeBits+bodySizeBits) | uint64(packetType)<<bodySizeBits | uint64(bodySize)
	endian.PutUint64(headerBytes[:], header)

	n, err := writer.Write(headerBytes[:])
	if err != nil {
		return err
	}

	if n != headerSize {
		return errWrittenSizeMismatch
	}

	return nil
}

func writePacketBody(writer io.Writer, body []byte) error {
	n, err := writer.Write(body)
	if err != nil {
		return err
	}

	if n != len(body) {
		return errWrittenSizeMismatch
	}

	return nil
}

func writePacket(writer io.Writer, packetType Type, body []byte) error {
	bodySize := int32(len(body))

	err := writePacketHeader(writer, packetType, bodySize)
	if err != nil {
		return err
	}

	if bodySize > 0 {
		err = writePacketBody(writer, body)
	}

	return nil
}
