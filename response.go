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
	SuccessReply = 0 // 成功的答复码
	ErrorReply   = 1 // 发生错误的答复码
)

// 从 reader 中读取数据并解析出响应内容。
func readResponseFrom(reader io.Reader) (reply byte, body []byte, err error) {

	// 读取指定字节数据
	header := make([]byte, headerLengthInProtocol)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return ErrorReply, nil, err
	}

	// 头部的第一个字节是协议版本号，如果版本号不一致很可能解析不成功，所以需要检查
	// 实际上这边可以做一个降级处理，就是尝试以响应的版本号去解析
	version := header[0]
	if version != ProtocolVersion {
		return ErrorReply, nil, errors.New("response " + ProtocolVersionMismatchErr.Error())
	}

	// 从头部解析出答复码还有响应体长度，同理，使用大端解析数字
	reply = header[1]
	header = header[2:]
	body = make([]byte, binary.BigEndian.Uint32(header))
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return ErrorReply, nil, err
	}
	return reply, body, nil
}

// 将响应写入到 writer。
func writeResponseTo(writer io.Writer, reply byte, body []byte) (int, error) {

	// 将响应体相关数据写入响应缓存区，并发送
	bodyLengthBytes := make([]byte, bodyLengthInProtocol)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(len(body)))

	response := make([]byte, 2, headerLengthInProtocol+len(body))
	response[0] = ProtocolVersion
	response[1] = reply
	response = append(response, bodyLengthBytes...)
	response = append(response, body...)
	return writer.Write(response)
}

// 向 writer 写入错误信息为 msg 的响应。
func writeErrorResponseTo(writer io.Writer, msg string) (int, error) {
	return writeResponseTo(writer, ErrorReply, []byte(msg))
}
