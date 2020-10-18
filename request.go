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

// 从 reader 中读取请求，并解析出命令和参数。
func readRequestFrom(reader io.Reader) (command byte, args [][]byte, err error) {

	// 读取头部，指定具体的大小，使用 ReadFull 读取满指定字节数据，如果数据还没传输过来，这个方法会进行等待
	header := make([]byte, headerLengthInProtocol)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return 0, nil, err
	}

	// 头部的第一个字节是协议版本号，拿出来判断协议版本号是否一致
	version := header[0]
	if version != ProtocolVersion {
		return 0, nil, ProtocolVersionMismatchErr
	}

	// 头部的第二个字节是命令，后面的四个字节是参数个数
	command = header[1]
	header = header[2:]

	// 所有的整数到字节数组的转换使用大端形式，所以这里使用 BigEndian 来将头部后四个字节转换为一个 uint32 数字
	argsLength := binary.BigEndian.Uint32(header)
	args = make([][]byte, argsLength)
	if argsLength > 0 {
		// 读取参数长度，同理使用大端处理，并一次性读取满参数
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

// 将请求写入到 writer 中。
func writeRequestTo(writer io.Writer, command byte, args [][]byte) (int, error) {

	// 创建一个缓存区，并将协议版本号、命令和参数个数等写入缓存区
	request := make([]byte, headerLengthInProtocol)
	request[0] = ProtocolVersion
	request[1] = command
	binary.BigEndian.PutUint32(request[2:], uint32(len(args)))

	if len(args) > 0 {
		// 将参数都添加到缓存区
		argLength := make([]byte, argLengthInProtocol)
		for _, arg := range args {
			binary.BigEndian.PutUint32(argLength, uint32(len(arg)))
			request = append(request, argLength...)
			request = append(request, arg...)
		}
	}
	return writer.Write(request)
}
