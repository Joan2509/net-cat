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
