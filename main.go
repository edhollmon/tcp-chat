package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/edhollmon/tcp-chat/server"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	addr := flag.String("addr", ":3000", "address to listen on (e.g. :3000, 0.0.0.0:8080)")
	flag.Parse()

	s := server.NewSimpleTCPServer(*addr)
	s.Start()

	<-quit
	fmt.Println("Shutting down...")
}
