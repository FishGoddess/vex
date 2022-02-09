// Copyright 2022 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2022/01/15 01:22:13

package vex

import (
	"bufio"
	"errors"
	"net"
)

type Client interface {
	Do(tag Tag, req []byte) (rsp []byte, err error)
	Close() error
}

type defaultClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(network string, address string) (Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &defaultClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (c *defaultClient) Do(tag Tag, req []byte) (rsp []byte, err error) {
	err = writeTo(c.writer, tag, req)
	if err != nil {
		return nil, err
	}

	err = c.writer.Flush()
	if err != nil {
		return nil, err
	}

	tag, body, err := readFrom(c.reader)
	if err != nil {
		return nil, err
	}

	if tag == errTag {
		return body, errors.New(string(body))
	}
	return body, nil
}

func (c *defaultClient) Close() error {
	err := c.writer.Flush()
	if err != nil {
		return err
	}
	return c.conn.Close()
}
