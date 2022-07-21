// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// config stores all configuration of client and server.
type config struct {
	// ConnTimeout is the timeout of a connection and any call will return an error if one connection has timeout.
	// See net.Conn's SetDeadline.
	ConnTimeout time.Duration

	// ReadBufferSize is the buffer size using in reading.
	// This value can be smaller if your reading data are often smaller.
	// This value can be bigger if your reading data are often bigger.
	// Notice: it applies to client and server.
	ReadBufferSize uint32

	// WriteBufferSize is the buffer size using in writing.
	// This value can be smaller if your writing data are often smaller.
	// This value can be bigger if your writing data are often bigger.
	// Notice: it applies to client and server.
	WriteBufferSize uint32

	// MaxConnected is the max-opened count of connections.
	MaxConnected uint64

	// EventHandler is a handler for handling events.
	EventHandler EventHandler
}

// newDefaultConfig returns a default config.
func newDefaultConfig() *config {
	return &config{
		ConnTimeout:     8 * time.Hour,
		ReadBufferSize:  4 * 1024 * 1024, // 4 MB
		WriteBufferSize: 4 * 1024 * 1024, // 4 MB
		MaxConnected:    4096,
		EventHandler:    NewDefaultEventHandler(""),
	}
}

// ApplyOptions applies opts to config.
func (c *config) ApplyOptions(opts []Option) *config {
	for _, opt := range opts {
		opt(c)
	}
	return c
}
