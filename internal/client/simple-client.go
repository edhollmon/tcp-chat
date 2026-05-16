package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type SimpleTCPClient struct {
	address string
	conn    net.Conn
}

func NewSimpleTCPClient(address string) *SimpleTCPClient {
	return &SimpleTCPClient{
		address: address,
	}
}

func (c *SimpleTCPClient) Start() {
	if err := c.connect(); err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer c.conn.Close()

	go c.readLoop()
	c.writeLoop()
}

func (c *SimpleTCPClient) connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *SimpleTCPClient) writeLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Fprintf(c.conn, "%s\n", line); err != nil {
			fmt.Println("Send error:", err)
			return
		}
	}
}

func (c *SimpleTCPClient) readLoop() {
	defer c.conn.Close()

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		msg := scanner.Bytes()
		fmt.Println(string(msg))
	}
}
