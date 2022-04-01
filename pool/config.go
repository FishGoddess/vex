// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

const (
	// FullStrategyBlock means pool.Get will block util pool has an idle client.
	FullStrategyBlock = 0

	// FullStrategyFailed means pool.Get will return an error if pool doesn't have an idle client.
	FullStrategyFailed = 1

	// FullStrategyNew means pool.Get will return a new client if pool doesn't have an idle client.
	FullStrategyNew = 2
)

// Option applies functions to config.
type Option func(c *config)

// WithMaxOpened sets maxOpened to config.
func WithMaxOpened(maxOpened int) Option {
	return func(c *config) {
		c.maxOpened = maxOpened
	}
}

// WithMaxIdle sets maxIdle to config.
func WithMaxIdle(maxIdle int) Option {
	return func(c *config) {
		c.maxIdle = maxIdle
	}
}

// WithFullStrategy sets strategy to config.
// See FullStrategy.
func WithFullStrategy(strategy FullStrategy) Option {
	return func(c *config) {
		c.fullStrategy = strategy
	}
}

// FullStrategy decides what pool will do when it's full.
type FullStrategy int8

// Block returns if fs is block full strategy.
func (fs FullStrategy) Block() bool {
	return fs == FullStrategyBlock
}

// Failed returns if fs is failed full strategy.
func (fs FullStrategy) Failed() bool {
	return fs == FullStrategyFailed
}

// New returns if fs is new full strategy.
func (fs FullStrategy) New() bool {
	return fs == FullStrategyNew
}

// config stores all configuration of Pool.
type config struct {
	// maxOpened is the max opened count of connections.
	maxOpened int

	// maxIdle is the max idle count of connections.
	maxIdle int

	// fullStrategy decides what pool will do when it's full.
	fullStrategy FullStrategy
}

// newDefaultConfig returns a default pool config.
func newDefaultConfig() config {
	return config{
		maxOpened:    64,
		maxIdle:      64,
		fullStrategy: FullStrategyBlock,
	}
}

// applyOptions applies opts to config.
func (c *config) applyOptions(opts []Option) {
	for _, opt := range opts {
		opt(c)
	}
}
