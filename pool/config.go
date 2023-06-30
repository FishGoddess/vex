// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

type Config struct {
	// MaxConnected is the max-opened count of connections.
	MaxConnected uint64

	// MaxIdle is the max-idle count of connections.
	MaxIdle uint64

	// BlockOnFull means getting clients from pool will block if connected is greater than max connected.
	BlockOnFull bool
}

func newDefaultConfig() *Config {
	return &Config{
		MaxConnected: 256,
		MaxIdle:      256,
		BlockOnFull:  true,
	}
}

func (c *Config) ApplyOptions(opts []Option) *Config {
	for _, opt := range opts {
		opt.ApplyTo(c)
	}

	return c
}
