// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

type Config struct {
	address string

	// name is a flag of server.
	name string

	// readTimeout is the timeout of reading from connection and any call will return an error if read timeout.
	// See net.Conn's SetReadDeadline.
	readTimeout time.Duration

	// writeTimeout is the timeout of writing to connection and any call will return an error if write timeout.
	// See net.Conn's SetWriteDeadline.
	writeTimeout time.Duration

	// closeTimeout is the timeout of closing a server.
	// Close may take a long time to wait all connections to be closed, so a timeout is necessary.
	closeTimeout time.Duration

	// readBufferSize is the buffer size used in reading.
	// This value can be smaller if your reading data are often smaller.
	// This value can be bigger if your reading data are often bigger.
	readBufferSize int

	// writeBufferSize is the buffer size used in writing.
	// This value can be smaller if your writing data are often smaller.
	// This value can be bigger if your writing data are often bigger.
	writeBufferSize int

	// maxConnections is the max number of connections.
	maxConnections int

	// beforeServingFunc is a function called before serving a server.
	beforeServingFunc func(address string)

	// afterServingFunc is a function called after serving a server.
	afterServingFunc func(address string, err error)

	// beforeHandlingFunc is a function called before handling a server.
	beforeHandlingFunc func(ctx *Context)

	// afterHandlingFunc is a function called after handling a server.
	afterHandlingFunc func(ctx *Context)

	// beforeClosingFunc is a function called before closing a server.
	beforeClosingFunc func(address string)

	// afterClosingFunc is a function called after closing a server.
	afterClosingFunc func(address string, err error)
}

func newClientConfig(address string) *Config {
	return &Config{
		address:         address,
		readTimeout:     10 * time.Minute,
		writeTimeout:    10 * time.Minute,
		readBufferSize:  64 * 1024, // 16KB
		writeBufferSize: 64 * 1024, // 16KB
	}
}

func newServerConfig(address string) *Config {
	return &Config{
		address:         address,
		name:            address,
		readTimeout:     10 * time.Minute,
		writeTimeout:    10 * time.Minute,
		closeTimeout:    time.Minute,
		readBufferSize:  16 * 1024, // 16KB
		writeBufferSize: 16 * 1024, // 16KB
		maxConnections:  4096,
	}
}

func (c *Config) ApplyOptions(opts []Option) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}

func (c *Config) beforeServing(address string) {
	if c.beforeServingFunc != nil {
		c.beforeServingFunc(address)
	}
}

func (c *Config) afterServing(address string, err error) {
	if c.afterServingFunc != nil {
		c.afterServingFunc(address, err)
	}
}

func (c *Config) beforeHandling(ctx *Context) {
	if c.beforeHandlingFunc != nil {
		c.beforeHandlingFunc(ctx)
	}
}

func (c *Config) afterHandling(ctx *Context) {
	if c.afterHandlingFunc != nil {
		c.afterHandlingFunc(ctx)
	}
}

func (c *Config) beforeClosing(address string) {
	if c.beforeClosingFunc != nil {
		c.beforeClosingFunc(address)
	}
}

func (c *Config) afterClosing(address string, err error) {
	if c.afterClosingFunc != nil {
		c.afterClosingFunc(address, err)
	}
}
