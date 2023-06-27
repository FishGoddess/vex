// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import "time"

// ClientOption applies functions to config.
type ClientOption func(conf *Config)

func (co ClientOption) ApplyTo(conf *Config) {
	co(conf)
}

// WithClientReadTimeout sets read timeout to config.
func WithClientReadTimeout(timeout time.Duration) ClientOption {
	return func(conf *Config) {
		conf.ReadTimeout = timeout
	}
}

// WithClientWriteTimeout sets write timeout to config.
func WithClientWriteTimeout(timeout time.Duration) ClientOption {
	return func(conf *Config) {
		conf.WriteTimeout = timeout
	}
}

// WithClientReadBufferSize sets read buffer size to config.
func WithClientReadBufferSize(bufferSize uint32) ClientOption {
	return func(conf *Config) {
		conf.ReadBufferSize = int(bufferSize)
	}
}

// WithClientWriteBufferSize sets write buffer size to config.
func WithClientWriteBufferSize(bufferSize uint32) ClientOption {
	return func(conf *Config) {
		conf.WriteBufferSize = int(bufferSize)
	}
}

// ServerOption applies functions to config.
type ServerOption func(conf *Config)

func (so ServerOption) ApplyTo(conf *Config) {
	so(conf)
}

// WithServerName sets name to config.
func WithServerName(name string) ServerOption {
	return func(conf *Config) {
		conf.Name = name
	}
}

// WithServerReadTimeout sets read timeout to config.
func WithServerReadTimeout(timeout time.Duration) ServerOption {
	return func(conf *Config) {
		conf.ReadTimeout = timeout
	}
}

// WithServerWriteTimeout sets write timeout to config.
func WithServerWriteTimeout(timeout time.Duration) ServerOption {
	return func(conf *Config) {
		conf.WriteTimeout = timeout
	}
}

// WithServerCloseTimeout sets close timeout to config.
func WithServerCloseTimeout(timeout time.Duration) ServerOption {
	return func(conf *Config) {
		conf.CloseTimeout = timeout
	}
}

// WithServerReadBufferSize sets read buffer size to config.
func WithServerReadBufferSize(bufferSize uint32) ServerOption {
	return func(conf *Config) {
		conf.ReadBufferSize = int(bufferSize)
	}
}

// WithServerWriteBufferSize sets write buffer size to config.
func WithServerWriteBufferSize(bufferSize uint32) ServerOption {
	return func(conf *Config) {
		conf.WriteBufferSize = int(bufferSize)
	}
}
