// Copyright 2023 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pack

import (
	"errors"

	"github.com/FishGoddess/vex"
)

// Send sends a packet to server and gets a packet responded from server.
func Send(client vex.Client, packetType PacketType, requestPacket []byte) (responsePacket []byte, err error) {
	err = writePacket(client, packetType, requestPacket)
	if err != nil {
		return nil, err
	}

	packetType, responsePacket, err = readPacket(client)
	if err != nil {
		return nil, err
	}

	if packetType == packetTypeErr {
		err = errors.New(string(responsePacket))
	}

	return responsePacket, err
}
