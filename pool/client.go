package pool

import (
	"errors"

	"github.com/FishGoddess/vex"
)

var (
	errClientClosed = errors.New("vex: client is closed")
)

// poolClient wraps client to a pool client.
type poolClient struct {
	pool   *Pool
	client vex.Client
}

// wrapClient wraps client to a pool client.
func wrapClient(pool *Pool, client vex.Client) vex.Client {
	return &poolClient{
		pool:   pool,
		client: client,
	}
}

// Send sends a packet with requestBody to server and returns responseBody responded from server.
func (pc *poolClient) Send(packetType vex.PacketType, requestBody []byte) (responseBody []byte, err error) {
	return pc.client.Send(packetType, requestBody)
}

// Close closes current client.
func (pc *poolClient) Close() error {
	pc.pool.putIdle(pc)
	return nil
}
