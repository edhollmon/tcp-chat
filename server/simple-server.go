package server

import (
	"fmt"
	"net"
)

type SimpleTCPServer struct {
	Listener    net.Listener
	Address     string
	Connections []net.Conn
}

type SimpleConnectionHandler = func(conn net.Conn)

func NewSimpleTCPServer(address string) *SimpleTCPServer {
	return &SimpleTCPServer{
		Address: address,
	}
}

func (s *SimpleTCPServer) Listen() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	fmt.Println("Listening on:", s.Address)
	s.Listener = listener
	return nil
}

func (s *SimpleTCPServer) HandleConnections(handler SimpleConnectionHandler) {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// TODO: Investigate if this is safe. May need Mutex?
		s.Connections = append(s.Connections, conn)
		go handler(conn)
	}
}
