// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"bufio"
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"

	packets "github.com/FishGoddess/vex/internal/packet"
)

// Client is the interface of vex client.
type Client interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
	Close() error
}

type client struct {
	conf *config

	conn     net.Conn
	id       atomic.Uint64
	inflight map[uint64]chan *packets.Packet
	done     chan struct{}

	once sync.Once
	lock sync.Mutex
}

// NewClient creates a client with address.
func NewClient(address string, opts ...Option) (Client, error) {
	conf := newConfig().apply(opts...)

	conn, err := net.DialTimeout("tcp", address, conf.dialTimeout)
	if err != nil {
		return nil, err
	}

	client := &client{
		conf:     conf,
		conn:     conn,
		inflight: make(map[uint64]chan *packets.Packet, 256),
		done:     make(chan struct{}),
	}

	go client.inflightLoop()
	return client, nil
}

func (c *client) inflightLoop() {
	reader := bufio.NewReader(c.conn)
	for {
		// Packet
		packet, err := packets.ReadPacket(reader)
		if err != nil {
			c.Close()
			return
		}

		// Inflight
		select {
		case <-c.done:
			return
		default:
			c.lock.Lock()
			ch := c.inflight[packet.ID()]
			c.lock.Unlock()

			if ch != nil {
				ch <- &packet
			}
		}
	}
}

// Send sends data and gets a new data.
// Returns an error if failed.
func (c *client) Send(ctx context.Context, data []byte) ([]byte, error) {
	// Inflight
	id := c.id.Add(1)
	ch := make(chan *packets.Packet, 1)

	c.lock.Lock()
	if c.inflight == nil {
		c.lock.Unlock()

		err := errors.New("vex: client is closed")
		return nil, err
	}

	c.inflight[id] = ch
	c.lock.Unlock()

	defer func() {
		c.lock.Lock()
		delete(c.inflight, id)
		c.lock.Unlock()
	}()

	// Packet
	packet := packets.New(id)
	packet.SetData(data)

	err := packets.WritePacket(c.conn, packet)
	if err != nil {
		return nil, err
	}

	select {
	case packet := <-ch:
		if packet == nil {
			err := errors.New("vex: client is closed")
			return nil, err
		}

		return packet.Data()
	case <-ctx.Done():
		err = ctx.Err()
		return nil, err
	case <-c.done:
		err := errors.New("vex: client is closed")
		return nil, err
	}
}

// Close closes the client and returns an error if failed.
func (c *client) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.conn.Close(); err != nil {
		return err
	}

	for _, ch := range c.inflight {
		ch <- nil
	}

	c.once.Do(func() { close(c.done) })
	c.id.Store(0)
	c.inflight = nil
	return nil
}
