package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

type SimpleTCPServer struct {
	mu       sync.RWMutex
	Listener net.Listener
	Address  string

	// Go Routine Trackers
	grWG sync.WaitGroup

	// Client Management
	clients map[uint64]*client
	nextcid uint64
}

type client struct {
	srv  *SimpleTCPServer
	cid  uint64
	conn net.Conn
}

func (c *client) readLoop() {
	defer c.srv.grWG.Done()
	defer c.conn.Close()
	defer func() {
		c.srv.mu.Lock()
		delete(c.srv.clients, c.cid)
		c.srv.mu.Unlock()
		fmt.Printf("Client %d disconnected\n", c.cid)
		c.srv.broadcast(fmt.Appendf(nil, "Client %d has left the chat\n", c.cid), c.cid)
	}()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		msg := scanner.Bytes()
		fmt.Println("Server received:", string(msg))
		c.srv.broadcast(fmt.Appendf(nil, "Client %d: %s\n", c.cid, msg), c.cid)
	}
}

func (s *SimpleTCPServer) broadcast(msg []byte, senderId uint64) {
	// Write to all other connections
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.clients {
		if c.cid == senderId {
			continue
		}
		c.conn.Write(msg)
	}
}

func NewSimpleTCPServer(address string) *SimpleTCPServer {
	return &SimpleTCPServer{
		mu:      sync.RWMutex{},
		Address: address,
		grWG:    sync.WaitGroup{},
		clients: make(map[uint64]*client),
		nextcid: 0,
	}
}

func (s *SimpleTCPServer) Start() {
	fmt.Println("Starting Simple TCP Server")
	defer fmt.Println("Server is ready")

	// TODO: Load Server Options here

	s.AcceptLoop()
}

func (s *SimpleTCPServer) AcceptLoop() {
	s.mu.Lock()
	l, err := s.getServerListener()
	if err != nil {
		s.mu.Unlock()
		fmt.Printf("Error listening on: %s", s.Address)
		fmt.Println(err)
		return
	}

	s.Listener = l

	s.grWG.Add(1)
	go s.acceptConnections(l)
	s.mu.Unlock()
}

func (s *SimpleTCPServer) getServerListener() (net.Listener, error) {
	l, err := net.Listen("tcp", s.Address)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (s *SimpleTCPServer) acceptConnections(l net.Listener) {
	defer s.grWG.Done()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Client connecting...")

		s.grWG.Add(1)
		go s.createClient(conn)
	}
}

func (s *SimpleTCPServer) createClient(conn net.Conn) {
	defer s.grWG.Done()

	s.mu.Lock()
	c := &client{
		cid:  atomic.AddUint64(&s.nextcid, 1),
		srv:  s,
		conn: conn,
	}

	s.clients[c.cid] = c
	s.mu.Unlock()

	fmt.Printf("Client %d connected\n", c.cid)
	s.broadcast(fmt.Appendf(nil, "Client %d has joined the chat\n", c.cid), c.cid)

	s.grWG.Add(1)
	go c.readLoop()
}

func (s *SimpleTCPServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.Listener.Close()
	for _, c := range s.clients {
		c.conn.Close()
	}
	s.mu.Unlock()

	done := make(chan struct{})
	go func() {
		s.grWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
