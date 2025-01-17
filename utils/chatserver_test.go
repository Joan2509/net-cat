package server

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestNewChatServer(t *testing.T) {
	// Remove test log file if it exists
	os.Remove("chat.log")

	server := newChatServer()
	defer server.logFile.Close()

	if server.clients == nil {
		t.Error("clients map was not initialized")
	}

	if server.messages == nil {
		t.Error("messages slice was not initialized")
	}

	if server.storedNames == nil {
		t.Error("storedNames map was not initialized")
	}

	if server.logFile == nil {
		t.Error("logFile was not initialized")
	}
}

func TestStart(t *testing.T) {
	// This is a basic smoke test for the Start function
	go func() {
		Start(8080)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
}
