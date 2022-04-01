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
func NewClient(network string, address string) (Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &defaultClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
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
	err := c.writer.Flush()
	if err != nil {
		return err
	}

	return c.conn.Close()
}
