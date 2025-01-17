package server

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

// mockConn implements net.Conn for testing
type mockConn struct {
	*bytes.Buffer
	readData  chan string
	writeData chan []byte
	closed    bool
}

func newMockConn() *mockConn {
	return &mockConn{
		Buffer:    new(bytes.Buffer),
		readData:  make(chan string, 100),
		writeData: make(chan []byte, 100),
	}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.closed {
		return 0, net.ErrClosed
	}
	data := <-m.readData
	copy(b, []byte(data))
	return len(data), nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	if m.closed {
		return 0, net.ErrClosed
	}
	m.writeData <- b
	return len(b), nil
}

func (m *mockConn) Close() error {
	m.closed = true
	close(m.readData)
	close(m.writeData)
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8081}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestHandleConnection(t *testing.T) {
	server := newChatServer()
	defer server.logFile.Close()

	conn := newMockConn()

	// Start a goroutine to handle the connection
	go server.handleConnection(conn)

	// Wait for welcome message and verify it
	var foundWelcome, foundPrompt bool
	timeout := time.After(time.Second)

	// Read messages until we find what we're looking for or timeout
	for {
		select {
		case data := <-conn.writeData:
			msg := string(data)
			if strings.Contains(msg, "Welcome to TCP-Chat!") {
				foundWelcome = true
			}
			if strings.Contains(msg, "[ENTER YOUR NAME]") {
				foundPrompt = true
			}
			if foundWelcome && foundPrompt {
				// Send the username
				go func() {
					conn.readData <- "TestUser\n"
				}()
				return
			}
		case <-timeout:
			if !foundWelcome {
				t.Error("Welcome message not received")
			}
			if !foundPrompt {
				t.Error("Name prompt not received")
			}
			return
		}
	}
}
