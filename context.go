// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"net"
	"time"
)

func setupConn(conf *Config, conn *net.TCPConn) error {
	now := time.Now()
	readDeadline := now.Add(conf.ReadTimeout)
	writeDeadline := now.Add(conf.WriteTimeout)

	if err := conn.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(writeDeadline); err != nil {
		return err
	}

	if err := conn.SetReadBuffer(conf.ReadBufferSize); err != nil {
		return err
	}

	if err := conn.SetWriteBuffer(conf.WriteBufferSize); err != nil {
		return err
	}

	return nil
}

type Context struct {
	parent context.Context
	cancel context.CancelFunc

	conn *net.TCPConn
}

func (c *Context) setup(conn *net.TCPConn) {
	c.parent, c.cancel = context.WithCancel(context.Background())
	c.conn = conn
}

func (c *Context) finish() (err error) {
	if err = c.conn.Close(); err != nil {
		return err
	}

	c.cancel()
	return nil
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *Context) Err() error {
	return c.parent.Err()
}

func (c *Context) Value(key any) any {
	return c.parent.Value(key)
}

func (c *Context) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *Context) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

func (c *Context) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Context) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
