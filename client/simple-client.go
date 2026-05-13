package client

import (
	"fmt"
	"net"
)

type SimpleTCPClient struct {
	Address string
	Conn    net.Conn
}

func NewSimpleTCPClient(address string) *SimpleTCPClient {
	return &SimpleTCPClient{
		Address: address,
	}
}

func (client *SimpleTCPClient) Connect() error {
	conn, err := net.Dial("tcp", client.Address)
	if err != nil {
		return err
	}

	client.Conn = conn
	return nil
}

func (client *SimpleTCPClient) Send(msg string) error {
	fmt.Println("Client attempting to send :", msg)
	_, err := client.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error writing message from client", err)
		return err
	}
	return nil
}
