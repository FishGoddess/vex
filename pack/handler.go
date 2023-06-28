// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"context"
	"errors"
	"io"
	"sync"
)

var (
	errPacketHandlerNotFound = errors.New("vex: packet handler not found")
)

// PacketHandler is a handler for handling packets.
// You will receive a byte slice of request and should return a byte slice or error if necessary.
type PacketHandler func(ctx context.Context, packetType PacketType, requestPacket []byte) (responsePacket []byte, err error)

type ServerHandler struct {
	handlers map[PacketType]PacketHandler
	lock     sync.RWMutex
}

func Handler() *ServerHandler {
	return &ServerHandler{
		handlers: make(map[PacketType]PacketHandler, 16),
	}
}

// RegisterPacketHandler registers a handler for packetType.
func (sh *ServerHandler) RegisterPacketHandler(packetType PacketType, handler PacketHandler) {
	sh.lock.Lock()
	sh.handlers[packetType] = handler
	sh.lock.Unlock()
}

func (sh *ServerHandler) writePacketOK(writer io.Writer, body []byte) {
	err := writePacket(writer, packetTypeOK, body)
	if err != nil {

	}
}

func (sh *ServerHandler) writePacketErr(writer io.Writer, err error) {
	err = writePacket(writer, packetTypeErr, []byte(err.Error()))
	if err != nil {

	}
}

func (sh *ServerHandler) Handle(ctx context.Context, reader io.Reader, writer io.Writer) {
	for {
		packetType, requestPacket, err := readPacket(reader)
		if err == io.EOF {
			return
		}

		if err != nil {
			sh.writePacketErr(writer, err)
			continue
		}

		sh.lock.RLock()
		handle, ok := sh.handlers[packetType]
		sh.lock.RUnlock()

		if !ok {
			sh.writePacketErr(writer, errPacketHandlerNotFound)
			continue
		}

		responsePacket, err := handle(ctx, packetType, requestPacket)
		if err != nil {
			sh.writePacketErr(writer, err)
			continue
		}

		sh.writePacketOK(writer, responsePacket)
	}
}
