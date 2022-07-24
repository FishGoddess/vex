// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"errors"
	"net"
)

// Client is the interface of vex client.
type Client interface {
	// Send sends a packet with requestBody to server and returns responseBody responded from server.
	Send(packetType PacketType, requestBody []byte) (responseBody []byte, err error)

	// Close closes current client.
	Close() error
}

// defaultClient is the default client implement which using one independent tcp connection.
type defaultClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewClient creates a new client to address with given network.
func NewClient(network string, address string, opts ...Option) (Client, error) {
	conn, err := dial(network, address)
	if err != nil {
		return nil, err
	}

	config := newDefaultConfig(network, address).ApplyOptions(opts)
	return &defaultClient{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, int(config.ReadBufferSize)),
		writer: bufio.NewWriterSize(conn, int(config.WriteBufferSize)),
	}, nil
}

// Send sends a packet with requestBody to server and returns responseBody responded from server.
func (c *defaultClient) Send(packetType PacketType, requestBody []byte) (responseBody []byte, err error) {
	err = writePacket(c.writer, packetType, requestBody)
	if err != nil {
		return nil, err
	}

	err = c.writer.Flush()
	if err != nil {
		return nil, err
	}

	packetType, responseBody, err = readPacket(c.reader)
	if err != nil {
		return nil, err
	}

	if packetType == packetTypeErr {
		return responseBody, errors.New(string(responseBody))
	}

	return responseBody, nil
}

// Close closes current client.
func (c *defaultClient) Close() error {
	if err := c.writer.Flush(); err != nil {
		return err
	}
	return c.conn.Close()
}
