package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func newChatServer() *ChatServer {
	logFile, err := os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	return &ChatServer{
		clients:  make(map[*Client]bool),
		messages: []string{},
		logFile:  logFile,
	}
}

func (s *ChatServer) start(port int) error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer s.listener.Close()

	log.Printf("Listening on port %d\n", port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *ChatServer) formatMessage(msg string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s]%s", timestamp, msg)
}

func (s *ChatServer) logMessage(msg string) {
	s.logFile.WriteString(msg + "\n")
	s.logFile.Sync()
}

func (s *ChatServer) broadcastMessage(msg string, sender *Client) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	for client := range s.clients {
		if client != sender {
			select {
			case client.messages <- msg:
			default:
				// If channel is full, remove client
				delete(s.clients, client)
			}
		}
	}
}

func (s *ChatServer) receiveMessages(client *Client) {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			continue // Skip empty messages
		}

		fullMsg := s.formatMessage(fmt.Sprintf("[%s]:%s", client.name, msg))
		s.broadcastMessage(fullMsg, client)
		s.logMessage(fullMsg)
	}

	// Client disconnected
	s.clientsMutex.Lock()
	delete(s.clients, client)
	s.clientsMutex.Unlock()

	disconnectMsg := s.formatMessage(fmt.Sprintf("%s has left our chat...", client.name))
	s.broadcastMessage(disconnectMsg, nil)
	s.logMessage(disconnectMsg)
}

func (s *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Check connection limit
	s.clientsMutex.Lock()
	if len(s.clients) >= maxConnections {
		conn.Write([]byte("Server is full. Try again later.\n"))
		s.clientsMutex.Unlock()
		return
	}
	s.clientsMutex.Unlock()

	// Send welcome message, Linux logo, and name prompt
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte(linuxLogo + "\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	// Get client name
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	clientName := scanner.Text()

	// Validate name
	if clientName == "" {
		conn.Write([]byte("Name cannot be empty. Disconnecting.\n"))
		return
	}

	client := &Client{
		name:     clientName,
		conn:     conn,
		messages: make(chan string, 100),
	}

	// Add client to the server's client list
	s.clientsMutex.Lock()
	s.clients[client] = true
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
	joinMsg := s.formatMessage(fmt.Sprintf("%s has joined our chat...", clientName))
	s.broadcastMessage(joinMsg, client)
	s.logMessage(joinMsg)

	// Handle client messages
	go s.receiveMessages(client)

	// Send messages to client
	for msg := range client.messages {
		conn.Write([]byte(msg + "\n"))
	}
}

func Start(port int) {
	server := newChatServer()
	log.Fatal(server.start(port))
}
