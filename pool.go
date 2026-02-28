// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"

	"github.com/FishGoddess/rego"
)

// Status is the status information of pool.
type Status rego.Status

type poolClient struct {
	pool *Pool

	client Client
}

// Send sends data and gets a new data.
// Returns an error if failed.
func (pc poolClient) Send(ctx context.Context, data []byte) ([]byte, error) {
	return pc.client.Send(ctx, data)
}

// Close returns the client back to the pool and returns an error if failed.
func (pc poolClient) Close() error {
	ctx := context.Background()
	return pc.pool.clients.Release(ctx, pc)
}

// DialFunc dials with context and returns the client.
// Returns an error if failed.
type DialFunc func(ctx context.Context) (Client, error)

// Pool is a pool for reusing clients.
// You should always use a pool instead of using a client in production.
type Pool struct {
	conf *config

	clients *rego.Pool[poolClient]
}

// NewPool returns a pool with limit and dial function.
// Dial function should return a new client as your way and an error if failed.
func NewPool(limit uint64, dial DialFunc, opts ...Option) *Pool {
	conf := newConfig().apply(opts...)
	pool := &Pool{conf: conf}

	acquire := func(ctx context.Context) (poolClient, error) {
		client, err := dial(ctx)
		if err != nil {
			return poolClient{}, err
		}

		pc := poolClient{pool: pool, client: client}
		return pc, nil
	}

	release := func(ctx context.Context, pc poolClient) error {
		return pc.client.Close()
	}

	pool.clients = rego.New(limit, acquire, release)
	return pool
}

// Get gets a client from pool and returns an error if failed.
func (p *Pool) Get(ctx context.Context) (Client, error) {
	return p.clients.Acquire(ctx)
}

// Status returns the status of pool.
func (p *Pool) Status() Status {
	status := p.clients.Status()
	return Status(status)
}

// Close closes the pool and releases all clients in it.
func (p *Pool) Close() error {
	ctx := context.Background()
	return p.clients.Close(ctx)
}
