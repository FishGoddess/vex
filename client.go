// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/17 17:47:21

package vex

import (
	"bufio"
	"errors"
	"io"
	"net"
)

type Client struct {
	conn   net.Conn
	reader io.Reader
}

func NewClient(network string, address string) (*Client, error) {

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) Do(command byte, args [][]byte) (body []byte, err error) {
	_, err = writeRequestTo(c.conn, command, args)
	if err != nil {
		return nil, err
	}

	reply, body, err := readResponseFrom(c.reader)
	if err != nil {
		return body, err
	}

	if reply == ErrorReply {
		return body, errors.New(string(body))
	}
	return body, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
