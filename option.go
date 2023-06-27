// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// ServerOption applies functions to server config.
type ServerOption func(conf *ServerConfig)

func (so ServerOption) ApplyTo(conf *ServerConfig) {
	so(conf)
}

// WithName sets name to server.
func WithName(name string) ServerOption {
	return func(conf *ServerConfig) {
		conf.Name = name
	}
}

// WithReadTimeout sets read timeout to server.
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(conf *ServerConfig) {
		conf.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets write timeout to server.
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(conf *ServerConfig) {
		conf.WriteTimeout = timeout
	}
}

// WithCloseTimeout sets close timeout to server.
func WithCloseTimeout(timeout time.Duration) ServerOption {
	return func(conf *ServerConfig) {
		conf.CloseTimeout = timeout
	}
}

// WithReadBufferSize sets read buffer size to config.
func WithReadBufferSize(bufferSize uint32) ServerOption {
	return func(conf *ServerConfig) {
		conf.ReadBufferSize = int(bufferSize)
	}
}

// WithWriteBufferSize sets write buffer size to server.
func WithWriteBufferSize(bufferSize uint32) ServerOption {
	return func(conf *ServerConfig) {
		conf.WriteBufferSize = int(bufferSize)
	}
}
