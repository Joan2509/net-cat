package main

import (
	"fmt"
	"os"
	"strconv"

	"net-cat/server"
)

func main() {
	defaultPort := 8989

	switch len(os.Args) {
	case 1:
		server.Start(defaultPort)
	case 2:
		p, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		server.Start(p)

	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
}
