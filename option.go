// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

type Option func(conf *Config)

func (o Option) ApplyTo(conf *Config) {
	o(conf)
}

// WithName sets name to config.
func WithName(name string) Option {
	return func(conf *Config) {
		conf.name = name
	}
}

// WithReadTimeout sets read timeout to config.
func WithReadTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.readTimeout = timeout
	}
}

// WithWriteTimeout sets write timeout to config.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.writeTimeout = timeout
	}
}

// WithConnectTimeout sets connect timeout to config.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.connectTimeout = timeout
	}
}

// WithCloseTimeout sets close timeout to config.
func WithCloseTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.closeTimeout = timeout
	}
}

// WithReadBufferSize sets read buffer size to config.
func WithReadBufferSize(bufferSize uint32) Option {
	return func(conf *Config) {
		conf.readBufferSize = int(bufferSize)
	}
}

// WithWriteBufferSize sets write buffer size to config.
func WithWriteBufferSize(bufferSize uint32) Option {
	return func(conf *Config) {
		conf.writeBufferSize = int(bufferSize)
	}
}

// WithMaxConnections sets max connections to config.
func WithMaxConnections(maxConnections uint32) Option {
	return func(conf *Config) {
		conf.maxConnections = int(maxConnections)
	}
}

// WithOnConnected sets on connected function to config.
func WithOnConnected(onConnected func(clientAddress string, serverAddress string)) Option {
	return func(conf *Config) {
		conf.onConnectedFunc = onConnected
	}
}

// WithOnDisconnected sets on disconnected function to config.
func WithOnDisconnected(onDisconnected func(clientAddress string, serverAddress string)) Option {
	return func(conf *Config) {
		conf.onDisconnectedFunc = onDisconnected
	}
}

// WithBeforeServing sets before serving function to config.
func WithBeforeServing(beforeServing func(address string)) Option {
	return func(conf *Config) {
		conf.beforeServingFunc = beforeServing
	}
}

// WithAfterServing sets after serving function to config.
func WithAfterServing(afterServing func(address string)) Option {
	return func(conf *Config) {
		conf.afterServingFunc = afterServing
	}
}

// WithBeforeHandling sets before handling function to config.
func WithBeforeHandling(beforeHandling func(ctx *Context)) Option {
	return func(conf *Config) {
		conf.beforeHandlingFunc = beforeHandling
	}
}

// WithAfterHandling sets after handling function to config.
func WithAfterHandling(afterHandling func(ctx *Context)) Option {
	return func(conf *Config) {
		conf.afterHandlingFunc = afterHandling
	}
}

// WithBeforeClosing sets before closing function to config.
func WithBeforeClosing(beforeClosing func(address string)) Option {
	return func(conf *Config) {
		conf.beforeClosingFunc = beforeClosing
	}
}

// WithAfterClosing sets after closing function to config.
func WithAfterClosing(afterClosing func(address string)) Option {
	return func(conf *Config) {
		conf.afterClosingFunc = afterClosing
	}
}
