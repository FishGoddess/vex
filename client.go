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
	conn *net.TCPConn
}

func NewClient(address string, readBufferSize int, writeBufferSize int) (Client, error) {
	resolved, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(network, nil, resolved)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
	}, nil
}

func (c *client) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *client) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *client) Close() error {
	return c.conn.Close()
}
