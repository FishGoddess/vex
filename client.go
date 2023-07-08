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

	conn        net.Conn
	connAddress string
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

func (c *client) connect() (err error) {
	defer func() {
		if err == nil {
			c.onConnected(c.connAddress, c.address)
			log.Debug("client %s has connected to %s", c.connAddress, c.address)
		}
	}()

	conn, err := net.DialTimeout(network, c.address, c.connectTimeout)
	if err != nil {
		return err
	}

	if err = setupConn(&c.Config, conn); err != nil {
		return err
	}

	c.conn = conn
	c.connAddress = c.conn.LocalAddr().String()

	return nil
}

// Read reads data to p.
// See io.Reader.
func (c *client) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

// Write writes p to data.
// See io.Writer.
func (c *client) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

// Close closes the client and returns an error if failed.
func (c *client) Close() (err error) {
	defer func() {
		if err == nil {
			c.onDisconnected(c.connAddress, c.address)
			log.Debug("client %s has disconnected from %s", c.connAddress, c.address)
		}
	}()

	return c.conn.Close()
}
