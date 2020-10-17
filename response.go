// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 15:48:20

package vex

import (
	"encoding/binary"
	"errors"
	"io"
)

// Response:
// version    reply    bodyLength    {body}
//  1byte     1byte      4byte      unknown

const (
	SuccessReply = 0
	ErrorReply   = 1
)

func readResponseFrom(reader io.Reader) (reply byte, body []byte, err error) {

	header := make([]byte, headerLengthInProtocol)
	_, err = io.ReadFull(reader, header)

	if err != nil {
		return ErrorReply, nil, err
	}

	version := header[0]
	if version != ProtocolVersion {
		return ErrorReply, nil, errors.New("response " + ProtocolVersionMismatchErr.Error())
	}

	reply = header[1]
	header = header[2:]
	body = make([]byte, binary.BigEndian.Uint32(header))
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return ErrorReply, nil, err
	}
	return reply, body, nil
}

func writeResponseTo(writer io.Writer, reply byte, body []byte) (int, error) {

	bodyLengthBytes := make([]byte, bodyLengthInProtocol)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(len(body)))

	response := make([]byte, 2, headerLengthInProtocol+len(body))
	response[0] = ProtocolVersion
	response[1] = reply
	response = append(response, bodyLengthBytes...)
	response = append(response, body...)
	return writer.Write(response)
}

func writeErrorResponseTo(writer io.Writer, msg string) (int, error) {
	return writeResponseTo(writer, ErrorReply, []byte(msg))
}
