package infrastructure

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// UnixSocketClient represents a client for Unix socket communication
type UnixSocketClient struct {
	socketPath string
	timeout    time.Duration
}

// NewUnixSocketClient creates a new Unix socket client
func NewUnixSocketClient(socketPath string, timeout time.Duration) *UnixSocketClient {
	return &UnixSocketClient{
		socketPath: socketPath,
		timeout:    timeout,
	}
}

// Get performs a GET request over Unix socket
func (c *UnixSocketClient) Get(path string) ([]byte, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, fmt.Errorf("connecting to unix socket %s: %w", c.socketPath, err)
	}
	defer conn.Close()

	// Set timeout
	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, fmt.Errorf("setting connection deadline: %w", err)
	}

	// Send HTTP request
	request := fmt.Sprintf("GET %s HTTP/1.0\r\nHost: localhost\r\n\r\n", path)
	if _, err := conn.Write([]byte(request)); err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	// Read response
	reader := bufio.NewReader(conn)

	// Parse HTTP response
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}