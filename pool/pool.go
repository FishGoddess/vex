// Copyright 2022 FishGoddess.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pool

import (
	"errors"
	"sync"

	"github.com/FishGoddess/vex"
)

var (
	errClientPoolClosed     = errors.New("vex: client pool is closed")
	errClientPoolFull       = errors.New("vex: client pool is full")
	errLimitStrategyUnknown = errors.New("vex: limit strategy is unknown")
)

// State stores all states of Pool.
type State struct {
	// Connected is the opened count of connections.
	Connected uint

	// Idle is the idle count of connections.
	Idle uint
}

// Pool is the pool of client.
type Pool struct {
	// config stores all configuration of Pool.
	config vex.Config

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
func NewPool(dial func() (vex.Client, error), opts ...vex.Option) *Pool {
	config := vex.NewDefaultConfig().ApplyOptions(opts)
	return &Pool{
		config:  *config,
		clients: make(chan *poolClient, config.MaxConnected),
		dial:    dial,
	}
}

// newClient returns a new Client.
func (cp *Pool) newClient() (vex.Client, error) {
	client, err := cp.dial()
	if err != nil {
		return nil, err
	}

	cp.state.Connected++
	return wrapClient(cp, client), nil
}

// putIdle stores a idle client to pool.
func (cp *Pool) putIdle(client *poolClient) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if cp.closed {
		client.client.Close()
		return
	}

	select {
	case cp.clients <- client:
		cp.state.Idle++
	default:
		client.client.Close()
	}
}

// getIdle gets an idle client from pool.
func (cp *Pool) getIdle() (*poolClient, bool) {
	select {
	case client := <-cp.clients:
		cp.state.Idle--
		return client, true
	default:
		return nil, false
	}
}

// getIdleBlocking gets an idle client from pool with blocking mode.
func (cp *Pool) getIdleBlocking() (*poolClient, bool) {
	client := <-cp.clients
	if client == nil {
		return nil, false
	}

	cp.lock.Lock()
	cp.state.Idle--
	cp.lock.Unlock()
	return client, true
}

// Get returns a client for use.
func (cp *Pool) Get() (vex.Client, error) {
	cp.lock.Lock()
	if cp.closed {
		cp.lock.Unlock()
		return nil, errClientPoolClosed
	}

	// Try to get an idle client.
	client, ok := cp.getIdle()
	if ok {
		cp.lock.Unlock()
		return client, nil
	}

	// Pool isn't full, returns a new client.
	if cp.state.Connected < cp.config.MaxConnected {
		defer cp.lock.Unlock()
		return cp.newClient()
	}

	// Pool is full:
	// 1. blocks util pool has an idle client.
	if cp.config.BlockOnLimit() {
		cp.lock.Unlock()

		client, ok = cp.getIdleBlocking()
		if ok {
			return client, nil
		}

		return nil, errClientClosed
	}

	// 2. returns an error.
	if cp.config.FailedOnLimit() {
		cp.lock.Unlock()
		return nil, errClientPoolFull
	}

	// 3. returns a new client.
	if cp.config.NewOnLimit() {
		defer cp.lock.Unlock()
		return cp.newClient()
	}

	cp.lock.Unlock()
	return nil, errLimitStrategyUnknown
}

// State returns all states of client pool.
func (cp *Pool) State() State {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	return cp.state
}

// Close closes pool and releases all resources.
func (cp *Pool) Close() error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if cp.closed {
		return nil
	}

	for i := uint(0); i < cp.state.Connected; i++ {
		client, ok := <-cp.clients
		if !ok {
			continue
		}

		if err := client.client.Close(); err != nil {
			return err
		}
	}

	cp.closed = true
	close(cp.clients)
	return nil
}
