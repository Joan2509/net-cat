package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type ChatServer struct {
	listener     net.Listener
	clients      map[*Client]bool
	clientsMutex sync.Mutex
	messages     []string
	logFile      *os.File
	storedNames  map[string]bool
}

func newChatServer() *ChatServer {
	logFile, err := os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return &ChatServer{
		clients:     make(map[*Client]bool),
		messages:    []string{},
		storedNames: make(map[string]bool),
		logFile:     logFile,
	}
}

func Start(port int) {
	server := newChatServer()

	// Start listening
	var err error
	server.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.listener.Close()

	log.Printf("Server is listening on port %d\n", port)

	// Accept connections
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

