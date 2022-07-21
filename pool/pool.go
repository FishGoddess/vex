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
	errClientPoolFull   = errors.New("vex: client pool is full")
	errClientPoolClosed = errors.New("vex: client pool is closed")
)

// State stores all states of Pool.
type State struct {
	// Connected is the opened count of connections.
	Connected uint64

	// Idle is the idle count of connections.
	Idle uint64
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
	}
}

// newClient returns a new Client.
func (cp *Pool) newClient() (vex.Client, error) {
	client, err := cp.dial()
	if err != nil {
		return nil, err
	}
	return wrapClient(cp, client), nil
}

// put adds an idle client to pool.
func (cp *Pool) put(client *poolClient) error {
	cp.lock.Lock()
	if cp.closed {
		cp.lock.Unlock()
		return client.client.Close()
	}

	if cp.state.Idle >= cp.config.MaxIdle {
		cp.lock.Unlock()
		return client.client.Close()
	}

	defer cp.lock.Unlock()
	select {
	case cp.clients <- client:
		cp.state.Idle++
		return nil
	default:
		return client.client.Close()
	}
}

// tryToGet tries to get an idle client from pool and return false if failed.
func (cp *Pool) tryToGet() (*poolClient, bool) {
	select {
	case client := <-cp.clients:
		return client, true
	default:
		return nil, false
	}
}

// waitToGet waits to get an idle client from pool.
func (cp *Pool) waitToGet() (*poolClient, bool) {
	client := <-cp.clients
	if client == nil {
		return nil, false
	}
	return client, true
}

// Get returns a client for use.
func (cp *Pool) Get() (vex.Client, error) {
	cp.lock.Lock()
	if cp.closed {
		cp.lock.Unlock()
		return nil, errClientPoolClosed
	}

	client, ok := cp.tryToGet()
	if ok {
		cp.state.Idle--
		cp.lock.Unlock()
		return client, nil
	}

	if cp.state.Connected < cp.config.MaxConnected {
		cp.state.Connected++
		cp.lock.Unlock()

		// Increase the connected and unlock before new client may cause the pool becomes full in advance.
		// So we should decrease the connected if new client failed.
		client, err := cp.newClient()
		if err != nil {
			cp.lock.Lock()
			cp.state.Connected--
			cp.lock.Unlock()
			return nil, err
		}

		return client, nil
	}

	if cp.config.BlockOnFull {
		cp.lock.Unlock()
		client, ok := cp.waitToGet()
		if ok {
			cp.lock.Lock()
			cp.state.Idle--
			cp.lock.Unlock()
			return client, nil
		}
		return nil, errClientPoolClosed
	}

	cp.lock.Unlock()
	return nil, errClientPoolFull
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

	for i := uint64(0); i < cp.state.Connected; i++ {
		client := <-cp.clients
		if client == nil {
			continue
		}

		if err := client.client.Close(); err != nil {
			return err
		}
	}

	close(cp.clients)
	cp.closed = true
	return nil
}
