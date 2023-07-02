// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"io"
	"net"

	"github.com/FishGoddess/vex/log"
)

type Client interface {
	io.ReadWriteCloser
}

type client struct {
	Config

	connection *Connection
}

// NewClient creates a new client connecting to address.
// Return an error if connect failed.
func NewClient(address string, opts ...Option) (Client, error) {
	client := &client{
		Config: *newClientConfig(address).ApplyOptions(opts),
	}

	return client, client.connect()
}

func (c *client) connect() error {
	resolved, err := net.ResolveTCPAddr(network, c.address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP(network, nil, resolved)
	if err != nil {
		return err
	}

	c.connection = newConnection(conn)
	if err = c.connection.setup(&c.Config); err != nil {
		return err
	}

	log.Debug("client %s has connected to %s", c.connection.LocalAddr(), c.connection.RemoteAddr())

	return nil
}

func (c *client) Read(p []byte) (n int, err error) {
	return c.connection.Read(p)
}

func (c *client) Write(p []byte) (n int, err error) {
	return c.connection.Write(p)
}

func (c *client) Close() error {
	return c.connection.close()
}
