// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 15:27:37

package vex

import (
	"encoding/binary"
	"io"
)

// Request:
// version    command    argsLength    {argLength    arg}
//  1byte      1byte       4byte          4byte    unknown

func readRequestFrom(reader io.Reader) (command byte, args [][]byte, err error) {

	header := make([]byte, headerLengthInProtocol)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return 0, nil, err
	}

	version := header[0]
	if version != ProtocolVersion {
		return 0, nil, ProtocolVersionMismatchErr
	}

	command = header[1]
	header = header[2:]
	argsLength := binary.BigEndian.Uint32(header)
	args = make([][]byte, argsLength)
	if argsLength > 0 {
		argLength := make([]byte, argLengthInProtocol)
		for i := uint32(0); i < argsLength; i++ {
			_, err = io.ReadFull(reader, argLength)
			if err != nil {
				return 0, nil, err
			}

			arg := make([]byte, binary.BigEndian.Uint32(argLength))
			_, err = io.ReadFull(reader, arg)
			if err != nil {
				return 0, nil, err
			}
			args[i] = arg
		}
	}
	return command, args, nil
}

func writeRequestTo(writer io.Writer, command byte, args [][]byte) (int, error) {

	request := make([]byte, headerLengthInProtocol)
	request[0] = ProtocolVersion
	request[1] = command
	binary.BigEndian.PutUint32(request[2:], uint32(len(args)))

	if len(args) > 0 {
		argLength := make([]byte, argLengthInProtocol)
		for _, arg := range args {
			binary.BigEndian.PutUint32(argLength, uint32(len(arg)))
			request = append(request, argLength...)
			request = append(request, arg...)
		}
	}
	return writer.Write(request)
}
