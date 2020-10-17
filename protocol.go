// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 15:15:57

package vex

import "errors"

// Request:
// version    command    argsLength    {argLength    arg}
//  1byte      1byte       4byte          4byte    unknown

// Response:
// version    reply    bodyLength    {body}
//  1byte     1byte      4byte      unknown

const (
	ProtocolVersion        = byte(1)
	headerLengthInProtocol = 6
	argsLengthInProtocol   = 4
	argLengthInProtocol    = 4
	bodyLengthInProtocol   = 4
)

var (
	ProtocolVersionMismatchErr = errors.New("protocol version between client and server doesn't match")
)
