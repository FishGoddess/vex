// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"io"
	"net"
)

type Client interface {
	io.ReadWriteCloser
}

type client struct {
	conn *Connection
}

func NewClient(address string, readBufferSize int, writeBufferSize int) (Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: newConnection(conn, readBufferSize, writeBufferSize),
	}, nil
}

func (c *client) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *client) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *client) Close() error {
	return c.conn.close()
}
