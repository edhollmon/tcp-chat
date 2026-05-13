package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":3000", "address to listen on (e.g. :3000, 0.0.0.0:8080)")
	flag.Parse()

	app := App{}
	app.Start(*addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down...")
}
