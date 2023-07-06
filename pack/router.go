// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/FishGoddess/vex"
	"github.com/FishGoddess/vex/log"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")
)

// PacketHandler is a handler for handling packets.
// You will receive a byte slice of request and should return a byte slice or error if necessary.
type PacketHandler func(ctx context.Context, packetType PacketType, requestPacket []byte) (responsePacket []byte, err error)

type Router struct {
	handlers map[PacketType]PacketHandler
	lock     sync.RWMutex
}

// NewRouter creates a router for registering some packet handlers.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[PacketType]PacketHandler, 16),
	}
}

// Register registers a packet handler to router.
func (r *Router) Register(packetType PacketType, handler PacketHandler) {
	r.lock.Lock()
	r.handlers[packetType] = handler
	r.lock.Unlock()
}

func (r *Router) writeStandardPacket(writer io.Writer, data []byte) {
	err := writePacket(writer, packetTypeStandard, data)
	if err != nil {
		log.Error(err, "write standard packet failed")
	}
}

func (r *Router) writeErrorPacket(writer io.Writer, err error) {
	err = writePacket(writer, packetTypeError, []byte(err.Error()))
	if err != nil {
		log.Error(err, "write error packet failed")
	}
}

// Handle handles context of vex, and you can pass it to a server.
func (r *Router) Handle(ctx *vex.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("router context of %s has done", ctx.RemoteAddr())
			return
		default:
		}

		packetType, data, err := readPacket(ctx)
		if err == io.EOF {
			return
		}

		if err != nil {
			r.writeErrorPacket(ctx, err)
			continue
		}

		r.lock.RLock()
		handle, ok := r.handlers[packetType]
		r.lock.RUnlock()

		if !ok {
			r.writeErrorPacket(ctx, errPacketHandlerNotFound)
			continue
		}

		data, err = handle(ctx, packetType, data)
		if err != nil {
			r.writeErrorPacket(ctx, err)
			continue
		}

		r.writeStandardPacket(ctx, data)
	}
}
