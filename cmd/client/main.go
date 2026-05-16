package main

import (
	"flag"
	"fmt"

	"github.com/edhollmon/tcp-chat/client"
)

func main() {
	addr := flag.String("addr", ":3000", "server address to connect to")
	flag.Parse()

	c := client.NewSimpleTCPClient(*addr)
	c.Start()

	fmt.Println("Connected to", *addr)
}
