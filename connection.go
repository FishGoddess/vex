// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"net"
	"time"
)

type Connection struct {
	conn *net.TCPConn
}

func newConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) setup(conf *Config) error {
	now := time.Now()
	readDeadline := now.Add(conf.ReadTimeout)
	writeDeadline := now.Add(conf.WriteTimeout)

	if err := c.conn.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	if err := c.conn.SetWriteDeadline(writeDeadline); err != nil {
		return err
	}

	if err := c.conn.SetReadBuffer(conf.ReadBufferSize); err != nil {
		return err
	}

	if err := c.conn.SetWriteBuffer(conf.WriteBufferSize); err != nil {
		return err
	}

	return nil
}

func (c *Connection) close() (err error) {
	return c.conn.Close()
}

func (c *Connection) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *Connection) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
