package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/edhollmon/tcp-chat/internal/server"
)

func main() {
	addr := flag.String("addr", ":3000", "address to listen on (e.g. :3000, 0.0.0.0:8080)")
	enableTrace := flag.Bool("trace", false, "enable runtime tracing to trace.out")
	flag.Parse()

	if *enableTrace {
		f, _ := os.Create("trace.out")
		defer f.Close()
		trace.Start(f)
		defer trace.Stop()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s := server.NewSimpleTCPServer(*addr)
	s.Start()

	<-quit
	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown timed out:", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped cleanly")
}
