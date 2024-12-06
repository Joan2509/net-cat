package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func Connect(host, port string) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Could not connect to server:", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Receive messages from server
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Send messages to server
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Text()
			fmt.Fprintf(conn, "%s\n", msg)
		}
	}()

	wg.Wait()
}
