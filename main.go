package main

import (
	"fmt"
	"os"

	"net-cat/utils"
)

func atoi(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			fmt.Println("Invalid character found")
			return -1
		}
	}
	return result
}

func main() {
	defaultPort := 8989

	switch len(os.Args) {
	case 1:
		server.Start(defaultPort)
	case 2:
		p := atoi(os.Args[1])
		if p < 1024 || p > 49151 {
			fmt.Println("[USAGE]: ./TCPChat $port")
			os.Exit(1)
		}
		server.Start(p)

	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
}
