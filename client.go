// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/13 21:45:22

package vex

import (
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(network string, address string) (*Client, error) {

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Do(command string, args [][]byte) ([]byte, error) {

	err := writeRequest(c.conn, &request{
		version: ProtocolVersion,
		command: command,
		args:    args,
	})

	if err != nil {
		return nil, err
	}
	return readResponse(c.conn)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
