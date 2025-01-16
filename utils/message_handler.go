package server

import (
	"bufio"
	"fmt"
	"time"
)

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

	// Add the message to the chat history
	s.messages = append(s.messages, msg)

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
		rawMsg := scanner.Text()
		if rawMsg == "" {
			continue // Skip empty messages
		}

		// Check if the message is a name change command
		if len(rawMsg) > 6 && rawMsg[:6] == "/name " {
			newName := rawMsg[6:] // Extract the new name
			if newName == "" {
				client.messages <- "Name cannot be empty. Try again."
				continue
			}

			// Broadcast the name change
			oldName := client.name
			client.name = newName
			notification := fmt.Sprintf("%s changed their name to %s.", oldName, newName)
			s.broadcastMessage(notification, nil)
			s.logMessage(notification)
			continue
		}

		fullMsg := s.formatMessage(fmt.Sprintf("[%s]:%s", client.name, rawMsg))

		// Clear the raw message from the sender's terminal
		client.messages <- "\033[1A\033[2K" // ANSI escape codes to move up and clear the line

		client.messages <- fullMsg
		s.broadcastMessage(fullMsg, client)
		s.logMessage(fullMsg)
	}

	// Client disconnected
	s.clientsMutex.Lock()
	delete(s.clients, client)
	delete(s.storedNames, client.name)
	s.clientsMutex.Unlock()

	disconnectMsg := fmt.Sprintf("%s has left our chat...", client.name)
	s.broadcastMessage(disconnectMsg, nil)
	s.logMessage(disconnectMsg)
}
