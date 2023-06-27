// Copyright 2023 FishGoddess. All rights reserved.
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
		conf.Name = name
	}
}

// WithReadTimeout sets read timeout to config.
func WithReadTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.ReadTimeout = timeout
	}
}

// WithWriteTimeout sets write timeout to config.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.WriteTimeout = timeout
	}
}

// WithCloseTimeout sets close timeout to config.
func WithCloseTimeout(timeout time.Duration) Option {
	return func(conf *Config) {
		conf.CloseTimeout = timeout
	}
}

// WithReadBufferSize sets read buffer size to config.
func WithReadBufferSize(bufferSize uint32) Option {
	return func(conf *Config) {
		conf.ReadBufferSize = int(bufferSize)
	}
}

// WithWriteBufferSize sets write buffer size to config.
func WithWriteBufferSize(bufferSize uint32) Option {
	return func(conf *Config) {
		conf.WriteBufferSize = int(bufferSize)
	}
}
