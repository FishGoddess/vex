// Copyright 2025 FishGoddess. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vex

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	packets "github.com/FishGoddess/vex/internal/packet"
)

func runTestServer() (string, func(), error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if errors.Is(err, net.ErrClosed) {
				break
			}

			if err != nil {
				continue
			}

			go func() {
				defer conn.Close()

				for {
					packet, err := packets.ReadPacket(conn)
					if err != nil {
						return
					}

					data, err := packet.Data()
					if err != nil {
						return
					}

					ii, err := strconv.Atoi(string(data))
					if err != nil {
						return
					}

					if ii%2 == 0 {
						err = errors.New(string(data))
						packet.SetError(err)
					}

					err = packets.WritePacket(conn, packet)
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	address := listener.Addr().String()
	done := func() { listener.Close() }
	return address, done, nil
}

// go test -v -cover -run=^TestNewClient$
func TestNewClient(t *testing.T) {
	ctx := context.Background()
	zeroClient := new(client)

	_, err := zeroClient.Send(ctx, nil)
	if err != errClientClosed {
		t.Fatalf("got %+v != want %+v", err, errClientClosed)
	}

	client, err := NewClient("")
	if err == nil {
		t.Fatal("new client returns a nil error")
	}

	address, done, err := runTestServer()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	client, err = NewClient(address)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
}

// go test -v -cover -run=^TestClientSend$
func TestClientSend(t *testing.T) {
	address, done, err := runTestServer()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	client, err := NewClient(address)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	ctx := context.Background()

	var group sync.WaitGroup
	for i := 1; i <= 100; i++ {
		ii := i

		group.Go(func() {
			data := []byte(strconv.Itoa(ii))

			gotData, err := client.Send(ctx, data)
			if ii%2 == 0 {
				if err == nil {
					t.Error("send returns a nil error")
				}

				got := err.Error()
				want := string(data)
				if got != want {
					t.Errorf("got %s != want %s", got, want)
				}

				return
			}

			if err != nil {
				t.Error(err)
			}

			got := string(gotData)
			want := string(data)
			if got != want {
				t.Errorf("got %s != want %s", got, want)
			}
		})
	}

	group.Wait()
}

// go test -v -cover -run=^TestClientSendTimeout$
func TestClientSendTimeout(t *testing.T) {
	address, done, err := runTestServer()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	client, err := NewClient(address)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var group sync.WaitGroup
	for i := 1; i <= 100; i++ {
		ii := i

		group.Go(func() {
			if ii%2 == 0 {
				time.Sleep(150 * time.Millisecond)
			}

			data := []byte(strconv.Itoa(ii))

			gotData, err := client.Send(ctx, data)
			if ii%2 == 0 {
				if err == nil {
					t.Error("send returns a nil error")
				}

				got := err.Error()
				want := context.DeadlineExceeded.Error()
				if got != want {
					t.Errorf("got %s != want %s", got, want)
				}

				return
			}

			if err != nil {
				t.Error(err)
			}

			got := string(gotData)
			want := string(data)
			if got != want {
				t.Errorf("got %s != want %s", got, want)
			}
		})
	}

	group.Wait()
}

// go test -v -cover -run=^TestClientClose$
func TestClientClose(t *testing.T) {
	address, done, err := runTestServer()
	if err != nil {
		t.Fatal(err)
	}

	defer done()

	cli, err := NewClient(address)
	if err != nil {
		t.Fatal(err)
	}

	if err = cli.Close(); err != nil {
		t.Fatal(err)
	}

	client := cli.(*client)

	select {
	case <-client.ctx.Done():
		t.Log("client context is done")
	default:
		t.Fatal("client context not done")
	}

	if client.inflight != nil {
		t.Fatal("client inflight not nil")
	}

	if client.inflightID != 0 {
		t.Fatal("client inflight not zero")
	}

	ctx := context.Background()

	_, err = cli.Send(ctx, nil)
	if err != errClientClosed {
		t.Fatalf("got %+v != want %+v", err, errClientClosed)
	}
}
