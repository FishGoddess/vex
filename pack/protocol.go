// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	magicBits    = 24
	typeBits     = 8
	dataSizeBits = 32

	maxMagic    = 1<<magicBits - 1
	maxType     = 1<<typeBits - 1
	maxDataSize = 1<<dataSizeBits - 1

	// Ha! Guess what this number means?
	magicNumber = 0xC638B

	headerSize = (magicBits + typeBits + dataSizeBits) / 8 // Bytes
)

const (
	packetTypeStandard PacketType = 0
	packetTypeError    PacketType = 1
)

var (
	Endian = binary.BigEndian
)

var (
	ErrWrongMagicNumber  = errors.New("vex: wrong magic number")
	ErrReadSizeMismatch  = errors.New("vex: read size != expected size")
	ErrWriteSizeMismatch = errors.New("vex: write size != expected size")
)

// PacketType is the type of packet.
type PacketType = byte

func readPacketHeader(reader io.Reader) (PacketType, int32, error) {
	var headerBytes [headerSize]byte

	n, err := reader.Read(headerBytes[:])
	if err != nil {
		return 0, 0, err
	}

	if n != headerSize {
		return 0, 0, ErrReadSizeMismatch
	}

	header := Endian.Uint64(headerBytes[:])
	magic := (header >> (typeBits + dataSizeBits)) & maxMagic

	if magic != magicNumber {
		return 0, 0, ErrWrongMagicNumber
	}

	packetType := PacketType((header >> dataSizeBits) & maxType)
	bodySize := int32(header & maxDataSize)
	return packetType, bodySize, nil
}

func readPacketBody(reader io.Reader, bodySize int32) ([]byte, error) {
	// Memory may exceed if body size is too big.
	body := make([]byte, bodySize, bodySize)

	n, err := io.ReadFull(reader, body)
	if err == io.ErrUnexpectedEOF {
		return nil, ErrReadSizeMismatch
	}

	if err != nil {
		return nil, err
	}

	if n != int(bodySize) {
		return nil, ErrReadSizeMismatch
	}

	return body, nil
}

func readPacket(reader io.Reader) (PacketType, []byte, error) {
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

func writePacketHeader(writer io.Writer, packetType PacketType, bodySize int32) error {
	var headerBytes [headerSize]byte
	var header = magicNumber<<(typeBits+dataSizeBits) | uint64(packetType)<<dataSizeBits | uint64(bodySize)
	Endian.PutUint64(headerBytes[:], header)

	n, err := writer.Write(headerBytes[:])
	if err != nil {
		return err
	}

	if n != headerSize {
		return ErrWriteSizeMismatch
	}

	return nil
}

func writePacketBody(writer io.Writer, body []byte) error {
	n, err := writer.Write(body)
	if err != nil {
		return err
	}

	if n != len(body) {
		return ErrWriteSizeMismatch
	}

	return nil
}

func writePacket(writer io.Writer, packetType PacketType, body []byte) error {
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
