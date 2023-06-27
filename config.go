// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// ServerConfig stores all configurations of server.
type ServerConfig struct {
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

// newServerConfig returns a new server config.
func newServerConfig(network string, address string) *ServerConfig {
	return &ServerConfig{
		address:         address,
		Name:            network + "/" + address,
		ReadTimeout:     10 * time.Minute,
		WriteTimeout:    10 * time.Minute,
		CloseTimeout:    time.Minute,
		ReadBufferSize:  64 * 1024, // 64KB
		WriteBufferSize: 64 * 1024, // 64KB
		MaxConnections:  4096,
	}
}

// ApplyOptions applies opts to a server config.
func (sc *ServerConfig) ApplyOptions(opts []ServerOption) *ServerConfig {
	for _, opt := range opts {
		opt.ApplyTo(sc)
	}

	return sc
}
