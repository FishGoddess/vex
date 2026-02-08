// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"errors"
	"fmt"
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
	sequence uint64
	inflight map[uint64]chan packets.Packet
	done     chan struct{}

	lock sync.Mutex
}

// NewClient creates a client with address.
func NewClient(address string, opts ...Option) (Client, error) {
	conf := newConfig().apply(opts...)

	conn, err := net.DialTimeout("tcp", address, conf.connectTimeout)
	if err != nil {
		return nil, err
	}

	client := &client{
		conf:     conf,
		conn:     conn,
		sequence: 0,
		inflight: make(map[uint64]chan packets.Packet, 256),
		done:     make(chan struct{}),
	}

	go client.readLoop()
	return client, nil
}

func (c *client) readLoop() {
	for {
		packet, err := packets.Decode(c.conn)
		if err != nil {
			return
		}

		select {
		case <-c.done:
			return
		default:
			c.dispatch(packet)
		}
	}
}

func (c *client) dispatch(packet packets.Packet) {
	c.lock.Lock()
	ch := c.inflight[packet.Sequence]
	c.lock.Unlock()

	if ch == nil {
		return
	}

	select {
	case ch <- packet:
	default:
	}
}

func (c *client) send(ctx context.Context, sequence uint64, data []byte, ch chan packets.Packet) ([]byte, error) {
	packet := packets.Packet{Magic: packets.Magic, Type: packets.PacketTypeRequest, Sequence: sequence}
	packet.With(data)

	err := packets.Encode(c.conn, packet)
	if err != nil {
		return nil, err
	}

	select {
	case packet, ok := <-ch:
		if !ok {
			err = errors.New("vex: inflight channel is closed")
			return nil, err
		}

		if packet.Type == packets.PacketTypeResponse {
			return packet.Data, nil
		}

		if packet.Type == packets.PacketTypeError {
			err = errors.New(string(packet.Data))
			return nil, err
		}

		err = fmt.Errorf("vex: packet type %v is wrong", packet.Type)
		return nil, err
	case <-ctx.Done():
		err = ctx.Err()
		return nil, err
	case <-c.done:
		err := errors.New("vex: client is closed")
		return nil, err
	}
}

// Send sends data and gets a new data.
// Returns an error if failed.
func (c *client) Send(ctx context.Context, data []byte) ([]byte, error) {
	sequence := atomic.AndUint64(&c.sequence, 1)
	ch := make(chan packets.Packet, 1)

	c.lock.Lock()
	if c.inflight == nil {
		c.lock.Unlock()

		err := errors.New("vex: inflight map is nil")
		return nil, err
	}

	c.inflight[sequence] = ch
	c.lock.Unlock()

	defer func() {
		c.lock.Lock()
		delete(c.inflight, sequence)
		c.lock.Unlock()
	}()

	return c.send(ctx, sequence, data, ch)
}

// Close closes the client and returns an error if failed.
func (c *client) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.conn.Close(); err != nil {
		return err
	}

	for _, ch := range c.inflight {
		close(ch)
	}

	c.sequence = 0
	c.inflight = nil
	return nil
}
