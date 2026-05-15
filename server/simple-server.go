package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type SimpleTCPServer struct {
	mu       sync.RWMutex
	Listener net.Listener
	Address  string

	// Go Routine Trackers
	grMu sync.RWMutex
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
	defer c.conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		msg := buf[:n]
		if n > 0 {
			fmt.Println("Server received:", string(msg))
			formatted := fmt.Sprintf("Client %d: %s", c.cid, msg)
			c.srv.broadcast([]byte(formatted), c.cid)
		}
		if err == io.EOF || err != nil {
			break
		}
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
		grMu:    sync.RWMutex{},
		grWG:    sync.WaitGroup{},
		clients: make(map[uint64]*client),
		nextcid: 1,
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

	go s.acceptConnections(l, func(conn net.Conn) { s.createClient(conn) })
	s.mu.Unlock()
}

func (s *SimpleTCPServer) getServerListener() (net.Listener, error) {
	l, err := net.Listen("tcp", s.Address)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (s *SimpleTCPServer) acceptConnections(l net.Listener, createFunc func(conn net.Conn)) {

	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO: Identify any errors where we may need to break out of infinite loop
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Client connecting...")

		if !s.startGoRoutine(func() {
			createFunc(conn)
			s.grWG.Done()
		}) {
			conn.Close()
		}
	}
}

func (s *SimpleTCPServer) createClient(conn net.Conn) *client {
	s.mu.Lock()
	c := &client{
		cid:  atomic.AddUint64(&s.nextcid, 1),
		srv:  s,
		conn: conn,
	}

	s.clients[c.cid] = c

	s.startGoRoutine(func() { c.readLoop() })

	s.mu.Unlock()

	return c
}

func (s *SimpleTCPServer) startGoRoutine(f func()) bool {
	s.grMu.Lock()
	defer s.grMu.Unlock()
	s.grWG.Add(1)
	go func() {
		f()
	}()
	return true
}
