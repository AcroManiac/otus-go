package gotelnet

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
)

type Client struct {
	host   string
	port   string
	dialer *net.Dialer
	conn   net.Conn
	input  io.Reader
	output io.Writer
}

func NewTelnetClient(host, port string, input io.Reader, output io.Writer) *Client {
	return &Client{
		host:   host,
		port:   port,
		dialer: nil,
		conn:   nil,
		input:  input,
		output: output,
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

func (c *Client) Send(ctx context.Context) (err error) {
	scanner := bufio.NewScanner(c.input)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()

			if _, err = c.conn.Write([]byte(fmt.Sprintf("%s\n", str))); err != nil {
				break OUTER
			}
		}
	}
	return
}

func (c *Client) Receive(ctx context.Context) (err error) {
	writer := bufio.NewWriter(c.output)
	scanner := bufio.NewScanner(c.conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			_, _ = writer.WriteString(scanner.Text())
			_ = writer.Flush()
		}
	}
	return
}

func (c *Client) Close() error {
	return c.conn.Close()
}
