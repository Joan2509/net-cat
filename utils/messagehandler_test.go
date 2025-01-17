package server

import (
	"strings"
	"testing"
	"time"
)

func TestFormatMessage(t *testing.T) {
	server := newChatServer()
	defer server.logFile.Close()

	msg := "test message"
	formatted := server.formatMessage(msg)

	if !strings.Contains(formatted, msg) {
		t.Errorf("Formatted message does not contain original message: %s", formatted)
	}

	if !strings.Contains(formatted, "[202") { // Year should start with 202x
		t.Errorf("Formatted message does not contain timestamp: %s", formatted)
	}
}

func TestBroadcastMessage(t *testing.T) {
	server := newChatServer()
	defer server.logFile.Close()

	// Create test clients
	client1 := &Client{
		name:     "User1",
		conn:     newMockConn(),
		messages: make(chan string, 100),
	}

	client2 := &Client{
		name:     "User2",
		conn:     newMockConn(),
		messages: make(chan string, 100),
	}

	// Add clients to server
	server.clientsMutex.Lock()
	server.clients[client1] = true
	server.clients[client2] = true
	server.clientsMutex.Unlock()

	// Test broadcast
	testMsg := "Test broadcast message"
	server.broadcastMessage(testMsg, client1)

	// Verify message was added to history
	if len(server.messages) != 1 {
		t.Error("Message was not added to history")
	}

	// Verify client2 received the message but client1 didn't
	select {
	case msg := <-client2.messages:
		if msg != testMsg {
			t.Errorf("Received wrong message: %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 did not receive the message")
	}

	select {
	case <-client1.messages:
		t.Error("Sender should not receive their own message")
	case <-time.After(100 * time.Millisecond):
		// This is expected
	}
}

func TestReceiveMessages(t *testing.T) {
	server := newChatServer()
	defer server.logFile.Close()

	mockConn := newMockConn()
	client := &Client{
		name:     "TestUser",
		conn:     mockConn,
		messages: make(chan string, 100),
	}

	// Add client to server
	server.clientsMutex.Lock()
	server.clients[client] = true
	server.storedNames[client.name] = true
	server.clientsMutex.Unlock()

	// Start receiving messages in a goroutine
	done := make(chan bool)
	go func() {
		server.receiveMessages(client)
		done <- true
	}()

	// Test name change command
	mockConn.readData <- "/name NewName\n"
	time.Sleep(100 * time.Millisecond)

	// Verify name change
	if client.name != "NewName" {
		t.Error("Name change command did not work")
	}

	// Test message sending
	mockConn.readData <- "Test message\n"
	time.Sleep(100 * time.Millisecond)

	// Verify message was broadcasted
	if len(server.messages) == 0 {
		t.Error("No messages were recorded")
	}

	// Simulate disconnect by closing the connection
	mockConn.Close()

	// Wait for disconnection handling
	select {
	case <-done:
		// Verify client removal
		server.clientsMutex.Lock()
		if _, exists := server.clients[client]; exists {
			t.Error("Client was not removed after disconnect")
		}
		server.clientsMutex.Unlock()
	case <-time.After(time.Second):
		t.Error("Timeout waiting for client disconnection")
	}
}
