package gotelnet

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	network = "tcp"
	host    = "127.0.0.1"
	port    = "4242"
)

func CreateServer(t *testing.T) (l net.Listener, err error) {
	l, err = net.Listen(network, fmt.Sprintf("%s:%s", host, port))
	require.NoError(t, err, "Should be no error while TCP creating server")
	return
}

func CreateClient(t *testing.T, r io.Reader, w io.Writer) (c *Client, err error) {
	c = NewTelnetClient(host, port, r, w)
	err = c.Connect(context.Background())
	require.NoError(t, err, "Should be no error while TCP server dialing")
	require.NotNil(t, c.conn, "Connection shouldn't be nil")
	return
}

func TestClient_Connect(t *testing.T) {
	// Create test TCP server
	s, _ := CreateServer(t)
	defer s.Close()

	// Connect test server
	go func() {
		_, _ = CreateClient(t, os.Stdin, os.Stdout)
	}()

	// Wait for client operations
	time.Sleep(time.Second)
}

func TestClient_Close(t *testing.T) {
	s, _ := CreateServer(t)
	defer s.Close()

	// Connect/close client
	go func() {
		client, _ := CreateClient(t, os.Stdin, os.Stdout)
		err := client.Close()
		require.NoError(t, err, "Should be no error while closing TCP server")
	}()

	// Wait for client operations
	time.Sleep(time.Second)
}

var sendString = "Answer to the Ultimate Question of Life, the Universe and Everything is 42"

func TestClient_Send(t *testing.T) {
	// Create test TCP server
	s, _ := CreateServer(t)
	defer s.Close()

	// Connect/send/close client
	go func() {
		r := strings.NewReader(sendString)
		client, _ := CreateClient(t, r, os.Stdout)

		err := client.Send(context.Background())
		require.NoError(t, err, "Should be no error while sending test data")

		err = client.Close()
		require.NoError(t, err, "Should be no error while closing TCP server")
	}()

	// Wait for client operations
	time.Sleep(time.Second)
}

func TestClient_Receive(t *testing.T) {
	// Create test TCP server
	s, _ := CreateServer(t)
	defer s.Close()

	// Connect/receive/close client
	go func() {
		client, _ := CreateClient(t, os.Stdin, os.Stdout)

		err := client.Receive(context.Background())
		require.NoError(t, err, "Should be no error while receiving test data")

		err = client.Close()
		require.NoError(t, err, "Should be no error while closing TCP server")
	}()

	// Send test data to client
	go func() {
		conn, err := s.Accept()
		require.NoError(t, err, "Should be no error while accepting TCP connection")
		require.NotNil(t, conn, "Connection shouldn't be nil")

		// Sending test data
		_, err = conn.Write([]byte(sendString))
		require.NoError(t, err, "Should be no error while writing test data")

		err = conn.Close()
		require.NoError(t, err, "Should be no error while closing connection")
	}()

	// Wait for server and client operations
	time.Sleep(time.Second)
}
