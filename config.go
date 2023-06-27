// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// Config stores all configurations of client and server.
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

// newClientConfig returns a new client config.
func newClientConfig(address string) *Config {
	return &Config{
		address:         address,
		ReadTimeout:     10 * time.Minute,
		WriteTimeout:    10 * time.Minute,
		ReadBufferSize:  16 * 1024, // 16KB
		WriteBufferSize: 16 * 1024, // 16KB
	}
}

// newServerConfig returns a new server config.
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

// ApplyClientOptions applies client options to config.
func (c *Config) ApplyClientOptions(opts []ClientOption) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}

// ApplyServerOptions applies server options to config.
func (c *Config) ApplyServerOptions(opts []ServerOption) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}
