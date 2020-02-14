package gotelnet

import (
	"context"
	"fmt"
	"net"
)

type Client struct {
	host   string
	port   string
	dialer *net.Dialer
	conn   net.Conn
}

func NewTelnetClient(host, port string) *Client {
	return &Client{
		host:   host,
		port:   port,
		dialer: nil,
		conn:   nil,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	c.dialer = &net.Dialer{}

	var err error
	c.conn, err = c.dialer.DialContext(ctx, "tcp",
		fmt.Sprintf("%s:%s", c.host, c.port))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Send() error {
	return nil
}

func (c *Client) Receive() error {
	return nil
}
