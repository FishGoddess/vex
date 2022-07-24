// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// Option applies functions to config.
type Option func(c *config)

// WithName sets name to config.
func WithName(name string) Option {
	return func(c *config) {
		c.Name = name
	}
}

// WithConnTimeout sets connection timeout to config.
func WithConnTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.ConnTimeout = timeout
	}
}

// WithCloseTimeout sets close timeout to config.
func WithCloseTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.CloseTimeout = timeout
	}
}

// WithReadBufferSize sets bufferSize to config.
func WithReadBufferSize(bufferSize uint32) Option {
	return func(c *config) {
		c.ReadBufferSize = bufferSize
	}
}

// WithWriteBufferSize sets bufferSize to config.
func WithWriteBufferSize(bufferSize uint32) Option {
	return func(c *config) {
		c.WriteBufferSize = bufferSize
	}
}

// WithMaxConnected sets maxConnected to config.
func WithMaxConnected(maxConnected uint64) Option {
	return func(c *config) {
		c.MaxConnected = maxConnected
	}
}

// WithEventListener sets event listener to config.
func WithEventListener(listener EventListener) Option {
	return func(c *config) {
		c.EventListener = listener
	}
}
