// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

type Option func(conf *Config)

func (o Option) ApplyTo(conf *Config) {
	o(conf)
}

// WithMaxConnected sets maxConnected to config.
func WithMaxConnected(maxConnected uint64) Option {
	return func(conf *Config) {
		conf.MaxConnected = maxConnected
	}
}

// WithMaxIdle sets maxIdle to config.
func WithMaxIdle(maxIdle uint64) Option {
	return func(conf *Config) {
		conf.MaxIdle = maxIdle
	}
}

// WithConnections sets maxConnected and maxIdle to config.
func WithConnections(connections uint64) Option {
	return func(conf *Config) {
		conf.MaxConnected = connections
		conf.MaxIdle = connections
	}
}

// WithNonBlockOnFull sets non-block on full to config.
func WithNonBlockOnFull() Option {
	return func(conf *Config) {
		conf.BlockOnFull = false
	}
}
