// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2022/01/15 00:23:09

package vex

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	versionSize = 1                                // 1 Byte
	tagSize     = 1                                // 1 Byte
	bodySize    = 4                                // 4 Byte
	headerSize  = versionSize + tagSize + bodySize // 6 Byte

	ProtocolVersion = 1 // v1
)

const (
	okTag  Tag = 0
	errTag Tag = 1
)

var (
	errProtocolMismatch = errors.New("vex: protocol between client and server doesn't match")
)

type Tag = byte

func readFrom(reader io.Reader) (tag Tag, body []byte, err error) {
	header := make([]byte, headerSize)

	_, err = reader.Read(header)
	if err != nil {
		return errTag, nil, err
	}

	if header[0] != ProtocolVersion {
		return errTag, nil, errProtocolMismatch
	}

	bodySize := binary.BigEndian.Uint32(header[versionSize+tagSize : headerSize])

	body = make([]byte, bodySize)
	_, err = reader.Read(body)
	if err != nil {
		return errTag, nil, err
	}
	return header[1], body, nil
}

func writeTo(writer io.Writer, tag Tag, body []byte) (err error) {
	header := make([]byte, headerSize)
	header[0] = ProtocolVersion
	header[1] = tag
	binary.BigEndian.PutUint32(header[versionSize+tagSize:headerSize], uint32(len(body)))

	_, err = writer.Write(header)
	if err != nil {
		return err
	}

	_, err = writer.Write(body)
	if err != nil {
		return err
	}
	return nil
}
