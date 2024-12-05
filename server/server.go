package server

import (
	"log"
	"os"
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
