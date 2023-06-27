// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"net"
)

type Connection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func newConnection(conn net.Conn, readBufferSize int, writeBufferSize int) *Connection {
	return &Connection{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, readBufferSize),
		writer: bufio.NewWriterSize(conn, writeBufferSize),
	}
}

func (c *Connection) close() error {
	return c.conn.Close()
}

func (c *Connection) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *Connection) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func (c *Connection) Flush() error {
	return c.writer.Flush()
}
