// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"context"
	"errors"
	"sync"

	"github.com/FishGoddess/vex"
)

var (
	errClientPoolFull   = errors.New("vex: client pool is full")
	errClientPoolClosed = errors.New("vex: client pool is closed")
)

// State stores all states of Pool.
type State struct {
	// Connected is the opened count of connections.
	Connected uint64 `json:"connected"`

	// Idle is the idle count of connections.
	Idle uint64 `json:"idle"`

	// Waiting is the waiting count of getting requests.
	Waiting uint64 `json:"waiting"`
}

// Pool is the pool of client.
type Pool struct {
	// config stores all configuration of Pool.
	config config

	// state stores all states of Pool.
	state State

	// clients stores all unused connections.
	clients chan *poolClient

	// dial is for creating a new Client.
	dial func() (vex.Client, error)

	closed bool
	lock   sync.Mutex
}

// NewPool returns a client pool storing some clients.
func NewPool(dial func() (vex.Client, error), opts ...Option) *Pool {
	config := newDefaultConfig().ApplyOptions(opts)
	return &Pool{
		config:  *config,
		clients: make(chan *poolClient, config.MaxConnected),
		dial:    dial,
		closed:  false,
	}
}

// newClient returns a new Client.
func (p *Pool) newClient() (vex.Client, error) {
	client, err := p.dial()
	if err != nil {
		return nil, err
	}
	return wrapClient(p, client), nil
}

// put adds an idle client to pool.
func (p *Pool) put(client *poolClient) error {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()
		return client.client.Close()
	}

	// Only waiting count < idle count will close idle client immediately.
	if p.state.Waiting < p.state.Idle && p.state.Idle >= p.config.MaxIdle {
		p.state.Connected--
		p.lock.Unlock()
		return client.client.Close()
	}

	defer p.lock.Unlock()
	select {
	case p.clients <- client:
		p.state.Idle++
		return nil
	default:
		return client.client.Close()
	}
}

// tryToGet tries to get an idle client from pool and return false if failed.
func (p *Pool) tryToGet() (*poolClient, bool) {
	select {
	case client := <-p.clients:
		return client, true
	default:
		return nil, false
	}
}

// waitToGet waits to get an idle client from pool.
// TODO Add ctx.Done() to select will cause a performance problem...
// The select won't call runtime.selectgo if only one case in it, and runtime.selectgo has 2 methods which will cause a performance problem:
//     sellock(scases, lockorder)
//     sg := acquireSudog()
// So we don't know what to do yet, but we think timeout mechanism won't be supported util we solved it.
func (p *Pool) waitToGet(ctx context.Context) (*poolClient, error) {
	select {
	case client := <-p.clients:
		if client == nil {
			return nil, errClientPoolClosed
		}
		return client, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Get returns a client for use.
func (p *Pool) Get(ctx context.Context) (vex.Client, error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()
		return nil, errClientPoolClosed
	}

	client, ok := p.tryToGet()
	if ok {
		p.state.Idle--
		p.lock.Unlock()
		if client == nil {
			return nil, errClientPoolClosed
		}
		return client, nil
	}

	if p.state.Connected < p.config.MaxConnected {
		p.state.Connected++
		p.lock.Unlock()

		// Increase the connected and unlock before new client may cause the pool becomes full in advance.
		// So we should decrease the connected if new client failed.
		client, err := p.newClient()
		if err != nil {
			p.lock.Lock()
			p.state.Connected--
			p.lock.Unlock()
			return nil, err
		}

		return client, nil
	}

	if p.config.BlockOnFull {
		p.state.Waiting++
		p.lock.Unlock()

		client, err := p.waitToGet(ctx)
		if err != nil {
			p.lock.Lock()
			p.state.Waiting--
			p.lock.Unlock()
			return nil, err
		}

		p.lock.Lock()
		p.state.Idle--
		p.state.Waiting--
		p.lock.Unlock()
		return client, nil
	}

	p.lock.Unlock()
	return nil, errClientPoolFull
}

// State returns all states of client pool.
func (p *Pool) State() State {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.state
}

// Close closes pool and releases all resources.
func (p *Pool) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.closed {
		return nil
	}

	for i := uint64(0); i < p.state.Connected; i++ {
		client := <-p.clients
		if client == nil {
			continue
		}

		if err := client.client.Close(); err != nil {
			return err
		}
	}

	close(p.clients)
	p.closed = true
	return nil
}
