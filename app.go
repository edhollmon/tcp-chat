package main

import (
	"fmt"
	"io"
	"net"

	"github.com/edhollmon/tcp-chat/server"
)

type App struct {
	TCPServer *server.SimpleTCPServer
}

func (app *App) Start(addr string) {
	s := server.NewSimpleTCPServer(addr)
	if err := s.Listen(); err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	app.TCPServer = s
	go s.HandleConnections(pingPongHandler)
}

func pingPongHandler(conn net.Conn) {
	fmt.Println("Handling incoming connection...")
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			fmt.Println("Server received:", string(buf[:n]))
		}
		if err == io.EOF || err != nil {
			break
		}
	}
}
