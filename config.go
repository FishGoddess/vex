// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

const (
	// limitStrategyBlock decides to block util connected less than max connected.
	limitStrategyBlock = 1

	// limitStrategyFailed decides to return an error if connected is greater than max connected.
	limitStrategyFailed = 2

	// limitStrategyNew decides to create a new connection if connected is greater than max connected.
	limitStrategyNew = 3
)

// LimitStrategy decides what it will do when connected is greater than max connected.
type LimitStrategy int8

// Config stores all configuration of Pool.
type Config struct {
	// MaxConnected is the max opened count of connections.
	MaxConnected uint64

	// MaxIdle is the max idle count of connections.
	MaxIdle uint64

	// LimitStrategy decides what it will do when connected is greater than max connected.
	LimitStrategy LimitStrategy

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

	// EventHandler is a handler for handling events.
	EventHandler EventHandler
}

// NewDefaultConfig returns a default config.
func NewDefaultConfig() *Config {
	return &Config{
		MaxConnected:    64,
		MaxIdle:         64,
		LimitStrategy:   limitStrategyBlock,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		EventHandler:    NewDefaultEventHandler(""),
	}
}

// ApplyOptions applies opts to config.
func (c *Config) ApplyOptions(opts []Option) *Config {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// BlockOnLimit returns if config is block limit strategy.
func (c *Config) BlockOnLimit() bool {
	return c.LimitStrategy == limitStrategyBlock
}

// FailedOnLimit returns if config is failed limit strategy.
func (c *Config) FailedOnLimit() bool {
	return c.LimitStrategy == limitStrategyFailed
}

// NewOnLimit returns if config is new limit strategy.
func (c *Config) NewOnLimit() bool {
	return c.LimitStrategy == limitStrategyNew
}
