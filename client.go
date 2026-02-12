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

	packets "github.com/FishGoddess/vex/internal/packet"
)

// Client is the interface of vex client.
type Client interface {
	Send(ctx context.Context, data []byte) ([]byte, error)
	Close() error
}

type client struct {
	conf *config

	ctx    context.Context
	cancel context.CancelFunc

	conn       net.Conn
	inflight   map[uint64]chan *packets.Packet
	inflightID uint64

	group sync.WaitGroup
	lock  sync.Mutex
}

// NewClient creates a client with address.
func NewClient(address string, opts ...Option) (Client, error) {
	conf := newConfig().apply(opts...)

	conn, err := net.DialTimeout("tcp", address, conf.dialTimeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	inflight := make(map[uint64]chan *packets.Packet, 256)

	client := new(client)
	client.conf = conf
	client.ctx = ctx
	client.cancel = cancel
	client.conn = conn
	client.inflight = inflight

	client.group.Go(client.inflightLoop)
	return client, nil
}

func (c *client) inflightPacket(packet *packets.Packet) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		c.lock.Lock()
		ch := c.inflight[packet.ID()]
		c.lock.Unlock()

		if ch != nil {
			ch <- packet
		}

		return nil
	}
}

func (c *client) inflightLoop() {
	reader := bufio.NewReader(c.conn)
	for {
		packet, err := packets.ReadPacket(reader)
		if err != nil {
			c.Close()
			return
		}

		if err = c.inflightPacket(&packet); err != nil {
			return
		}
	}
}

func (c *client) handleData(data []byte) (packet packets.Packet, packetCh chan *packets.Packet, done func(), err error) {
	c.lock.Lock()
	if c.inflight == nil {
		c.lock.Unlock()

		return packet, nil, nil, errors.New("vex: client is closed")
	}

	c.inflightID++
	inflightID := c.inflightID

	packetCh = make(chan *packets.Packet, 1)
	c.inflight[inflightID] = packetCh
	c.lock.Unlock()

	done = func() {
		c.lock.Lock()
		delete(c.inflight, inflightID)
		c.lock.Unlock()
	}

	packet = packets.New(inflightID)
	packet.SetData(data)
	return packet, packetCh, done, nil
}

func (c *client) waitData(ctx context.Context, packetCh chan *packets.Packet) ([]byte, error) {
	select {
	case packet := <-packetCh:
		return packet.Data()
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, errors.New("vex: client is closed")
	}
}

// Send sends data and gets a new data.
// Returns an error if failed.
func (c *client) Send(ctx context.Context, data []byte) ([]byte, error) {
	packet, packetCh, done, err := c.handleData(data)
	if err != nil {
		return nil, err
	}

	defer done()

	err = packets.WritePacket(c.conn, packet)
	if err != nil {
		return nil, err
	}

	return c.waitData(ctx, packetCh)
}

// Close closes the client and returns an error if failed.
func (c *client) Close() error {
	c.lock.Lock()
	if err := c.conn.Close(); err != nil {
		c.lock.Unlock()

		return err
	}

	c.cancel()
	c.inflight = nil
	c.inflightID = 0
	c.lock.Unlock()
	c.group.Wait()
	return nil
}
