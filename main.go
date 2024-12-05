package main

import (
	"fmt"
	"os"
	"strconv"

	"net-cat/client"
	"net-cat/server"
)

func main() {
	// Default port
	defaultPort := 8989

	// Check argument count and mode
	switch len(os.Args) {
	case 1:
		fmt.Printf("Listening on the port :%d\n", defaultPort)
		server.Start(defaultPort)
	case 2:
		// Check if argument is a valid port for server
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		fmt.Printf("Listening on the port :%d\n", p)
		server.Start(p)
	case 3:
		// Client mode with host and port
		host := os.Args[1]
		port := os.Args[2]
		client.Connect(host, port)
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		fmt.Println("Or: ./TCPChat [host] [port]")
		os.Exit(1)
	}
}
