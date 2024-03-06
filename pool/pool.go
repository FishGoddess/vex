// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"errors"

	"github.com/FishGoddess/rego"
	"github.com/FishGoddess/vex"
)

var (
	ErrPoolIsFull   = errors.New("vex: pool is full")
	ErrPoolIsClosed = errors.New("vex: pool is closed")
)

// DialFunc is a function dials to somewhere and returns a client and error if failed.
type DialFunc func() (vex.Client, error)

// Dial returns a function which dials to address with opts.
// It's a convenient way used in creating a pool.
func Dial(address string, opts ...vex.Option) DialFunc {
	return func() (vex.Client, error) {
		return vex.NewClient(address, opts...)
	}
}

type Status struct {
	// Limit is the limit of connected clients.
	Limit uint64 `json:"limit"`

	// Connected is the count of connected clients.
	Connected uint64 `json:"connected"`

	// Idle is the count of idle clients.
	Idle uint64 `json:"idle"`

	// Waiting is the count of waiting for a client.
	Waiting uint64 `json:"waiting"`
}

type Pool struct {
	clients *rego.Pool[vex.Client]
}

func New(dial DialFunc, opts ...Option) *Pool {
	pool := new(Pool)

	acquire := func() (vex.Client, error) {
		client, err := dial()
		if err != nil {
			return nil, err
		}

		return newPoolClient(pool, client), nil
	}

	release := func(client vex.Client) error {
		if pclient, ok := client.(*poolClient); ok {
			return pclient.close()
		}

		return client.Close()
	}

	pool.clients = rego.New(acquire, release, newRegoOptions(opts)...)

	return pool
}

func newRegoOptions(opts []Option) []rego.Option {
	conf := newConfig()

	for _, opt := range opts {
		opt.ApplyTo(&conf)
	}

	var result []rego.Option
	if conf.limit > 0 {
		result = append(result, rego.WithLimit(conf.limit))
	}

	if conf.fastFailed {
		result = append(result, rego.WithFastFailed())
	}

	result = append(result, rego.WithPoolFullErr(func(ctx context.Context) error {
		return ErrPoolIsFull
	}))

	result = append(result, rego.WithPoolClosedErr(func(ctx context.Context) error {
		return ErrPoolIsClosed
	}))

	return result
}

func (p *Pool) put(client vex.Client) error {
	return p.clients.Put(client)
}

func (p *Pool) Get(ctx context.Context) (vex.Client, error) {
	return p.clients.Take(ctx)
}

// Status returns the status of the pool.
func (p *Pool) Status() Status {
	status := p.clients.Status()

	return Status{
		Limit:     status.Limit,
		Connected: status.Acquired,
		Idle:      status.Idle,
		Waiting:   status.Waiting,
	}
}

// Close closes pool and releases all resources.
func (p *Pool) Close() error {
	return p.clients.Close()
}
