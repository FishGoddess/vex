// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"log/slog"
	"time"
)

// Logger is for logging some messages in different levels.
// You can just use log/slog package which is one implement of it.
type Logger interface {
	Debug(msg string, kvs ...any)
	Info(msg string, kvs ...any)
	Error(msg string, kvs ...any)
}

type config struct {
	logger         Logger
	flushInterval  time.Duration
	connectTimeout time.Duration
}

func newConfig() *config {
	conf := &config{
		logger:         slog.Default(),
		flushInterval:  time.Second,
		connectTimeout: 3 * time.Second,
	}

	return conf
}

func (c *config) apply(opts ...Option) *config {
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Option configures the config for client or server.
type Option func(c *config)

// WithLogger sets the logger to config.
func WithLogger(logger Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

// WithFlushInterval sets the flush interval to config.
func WithFlushInterval(interval time.Duration) Option {
	return func(c *config) {
		c.flushInterval = interval
	}
}

// WithConnectTimeout sets the connect timeout to config.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.connectTimeout = timeout
	}
}
