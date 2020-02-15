package gotelnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
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

func (c *Client) Send(ctx context.Context) (err error) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				err = errors.New("cancel scanning user input")
				break OUTER
			}
			str := scanner.Text()
			//log.Printf("To server %v\n", str)

			if _, err = c.conn.Write([]byte(fmt.Sprintf("%s\n", str))); err != nil {
				break OUTER
			}
		}
	}
	//log.Printf("Finished Send")
	return
}

func (c *Client) Receive(ctx context.Context) (err error) {
	writer := bufio.NewWriter(os.Stdout)
	scanner := bufio.NewScanner(c.conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				err = errors.New("cancel scanning network connection")
				break OUTER
			}
			writer.WriteString(scanner.Text())
			writer.Flush()
		}
	}
	//log.Printf("Finished Receive")
	return
}

func (c *Client) Close() error {
	return c.conn.Close()
}
