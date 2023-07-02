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

	conn *net.TCPConn
}

// NewClient creates a new client connecting to address.
// Return an error if connect failed.
func NewClient(address string, opts ...Option) (Client, error) {
	conf := newClientConfig(address).ApplyOptions(opts)

	client := &client{
		Config: *conf,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
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

	if err = setupConn(&c.Config, conn); err != nil {
		return err
	}

	c.conn = conn
	log.Debug("client %s has connected to %s", c.conn.LocalAddr(), c.conn.RemoteAddr())

	return nil
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
