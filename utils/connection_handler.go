package server

import (
	"bufio"
	"fmt"
	"net"
)

func (s *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	s.clientsMutex.Lock()
	if len(s.clients) > 10 {
		conn.Write([]byte("Server is full. Try again later.\n"))
		s.clientsMutex.Unlock()
		return
	}
	s.clientsMutex.Unlock()

	// Send welcome message and Linux logo
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte(linuxLogo + "\n"))

	// Ask for a unique name
	scanner := bufio.NewScanner(conn)
	var clientName string
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		scanner.Scan()
		clientName = scanner.Text()

		s.clientsMutex.Lock()
		_, nameInUse := s.storedNames[clientName]
		s.clientsMutex.Unlock()

		if clientName == "" {
			conn.Write([]byte("Name cannot be empty. Please try again.\n"))
		} else if nameInUse {
			conn.Write([]byte("Name is already in use. Please choose another name.\n"))
		} else {
			break
		}
	}

	// Add client to the server's client list
	client := &Client{
		name:     clientName,
		conn:     conn,
		messages: make(chan string, 100),
	}

	s.clientsMutex.Lock()
	s.clients[client] = true
	s.storedNames[clientName] = true // Register the client's name
	s.clientsMutex.Unlock()

	// Send previous messages to the new client
	conn.Write([]byte("----- Chat History -----\n"))
	s.clientsMutex.Lock()
	for _, msg := range s.messages {
		conn.Write([]byte(msg + "\n"))
	}
	s.clientsMutex.Unlock()
	conn.Write([]byte("------------------------\n"))

	// Broadcast new client join
	joinMsg := fmt.Sprintf("%s has joined our chat...", clientName)
	s.broadcastMessage(joinMsg, client)
	s.logMessage(joinMsg)

	// Handle client messages
	go s.receiveMessages(client)

	// Send messages to client
	for msg := range client.messages {
		conn.Write([]byte(msg + "\n"))
	}

	// On disconnect, clean up
	s.clientsMutex.Lock()
	delete(s.clients, client)
	delete(s.storedNames, clientName) // Remove the client's name
	s.clientsMutex.Unlock()

	disconnectMsg := fmt.Sprintf("%s has left our chat...", clientName)
	s.broadcastMessage(disconnectMsg, nil)
	s.logMessage(disconnectMsg)
}
