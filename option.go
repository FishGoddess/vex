// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// Option applies functions to config.
type Option func(c *Config)

// WithMaxConnected sets maxConnected to config.
func WithMaxConnected(maxConnected uint) Option {
	return func(c *Config) {
		c.MaxConnected = maxConnected
	}
}

// WithMaxIdle sets maxIdle to config.
func WithMaxIdle(maxIdle uint) Option {
	return func(c *Config) {
		c.MaxIdle = maxIdle
	}
}

// WithBlockOnLimit sets block limit strategy to config.
func WithBlockOnLimit() Option {
	return func(c *Config) {
		c.LimitStrategy = limitStrategyBlock
	}
}

// WithFailedOnLimit sets failed limit strategy to config.
func WithFailedOnLimit() Option {
	return func(c *Config) {
		c.LimitStrategy = limitStrategyFailed
	}
}

// WithNewOnLimit sets new limit strategy to config.
func WithNewOnLimit() Option {
	return func(c *Config) {
		c.LimitStrategy = limitStrategyNew
	}
}

// WithReadBufferSize sets bufferSize to config.
func WithReadBufferSize(bufferSize uint32) Option {
	return func(c *Config) {
		c.ReadBufferSize = bufferSize
	}
}

// WithWriteBufferSize sets bufferSize to config.
func WithWriteBufferSize(bufferSize uint32) Option {
	return func(c *Config) {
		c.WriteBufferSize = bufferSize
	}
}

// WithEventHandler sets handler to config.
func WithEventHandler(handler EventHandler) Option {
	return func(c *Config) {
		c.EventHandler = handler
	}
}

// WithConnTimeout sets timeout to config.
func WithConnTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ConnTimeout = timeout
	}
}
