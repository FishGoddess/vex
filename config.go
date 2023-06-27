// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

type Config struct {
	address string

	// Name is a flag of server.
	Name string

	// ReadTimeout is the timeout of reading from connection and any call will return an error if read timeout.
	// See net.Conn's SetReadDeadline.
	ReadTimeout time.Duration

	// WriteTimeout is the timeout of writing to connection and any call will return an error if write timeout.
	// See net.Conn's SetWriteDeadline.
	WriteTimeout time.Duration

	// CloseTimeout is the timeout of closing a server.
	// Close may take a long time to wait all connections to be closed, so a timeout is necessary.
	CloseTimeout time.Duration

	// ReadBufferSize is the buffer size used in reading.
	// This value can be smaller if your reading data are often smaller.
	// This value can be bigger if your reading data are often bigger.
	ReadBufferSize int

	// WriteBufferSize is the buffer size used in writing.
	// This value can be smaller if your writing data are often smaller.
	// This value can be bigger if your writing data are often bigger.
	WriteBufferSize int

	// MaxConnections is the max number of connections.
	MaxConnections int
}

func newClientConfig(address string) *Config {
	return &Config{
		address:         address,
		ReadTimeout:     10 * time.Minute,
		WriteTimeout:    10 * time.Minute,
		ReadBufferSize:  64 * 1024, // 16KB
		WriteBufferSize: 64 * 1024, // 16KB
	}
}

func newServerConfig(address string) *Config {
	return &Config{
		address:         address,
		Name:            address,
		ReadTimeout:     10 * time.Minute,
		WriteTimeout:    10 * time.Minute,
		CloseTimeout:    time.Minute,
		ReadBufferSize:  16 * 1024, // 16KB
		WriteBufferSize: 16 * 1024, // 16KB
		MaxConnections:  4096,
	}
}

func (c *Config) ApplyOptions(opts []Option) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}
