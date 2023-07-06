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
	readDeadline := now.Add(conf.readTimeout)
	writeDeadline := now.Add(conf.writeTimeout)

	if err := conn.SetReadDeadline(readDeadline); err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(writeDeadline); err != nil {
		return err
	}

	if err := conn.SetReadBuffer(conf.readBufferSize); err != nil {
		return err
	}

	if err := conn.SetWriteBuffer(conf.writeBufferSize); err != nil {
		return err
	}

	return nil
}

// Context connects client and server which can be read and written.
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

// Deadline returns the time when context has done.
// See context.Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

// Done returns a channel that's closed when context has done.
// See context.Context.
func (c *Context) Done() <-chan struct{} {
	return c.parent.Done()
}

// Err returns the underlying error.
// See context.Context.
func (c *Context) Err() error {
	return c.parent.Err()
}

// Value returns the value associated with this context for key, or nil if no value is associated with key.
// See context.Context.
func (c *Context) Value(key any) any {
	return c.parent.Value(key)
}

// Read reads data to p.
// See io.Reader.
func (c *Context) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

// Write writes p to data.
// See io.Writer.
func (c *Context) Write(p []byte) (n int, err error) {
	return c.conn.Write(p)
}

// LocalAddr returns the local network address.
func (c *Context) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Context) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
