package network

import (
	"context"
	"fmt"
	"github.com/hello-world/common/logger"
	"net"
	"sync"
)

// TCPServer represents the TCP server instance
type TCPServer struct {
	listener   net.Listener
	workerPool *WorkerPool
	clients    sync.Map
	quit       chan bool
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(address string, numWorkers int) (*TCPServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to create listener: %v", err)
		return nil, fmt.Errorf("failed to create listener: %v", err)
	}

	server := &TCPServer{
		listener: listener,
		quit:     make(chan bool),
	}

	// Initialize worker pool with a custom handler
	server.workerPool = NewWorkerPool(numWorkers, 1000, server.handleTask)
	logger.Info("Created TCP server with %d workers", numWorkers)

	return server, nil
}

// Start begins accepting connections
func (s *TCPServer) Start() {
	logger.Info("Starting TCP server")
	go s.acceptConnections()
}

// acceptConnections handles incoming connections
func (s *TCPServer) acceptConnections() {
	for {
		select {
		case <-s.quit:
			logger.Info("Stopping connection acceptance")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				logger.Error("Failed to accept connection: %v", err)
				continue
			}
			logger.Debug("Accepted new connection from %s", conn.RemoteAddr())
			go s.handleConnection(conn)
		}
	}
}

// handleConnection processes individual connections
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientID := conn.RemoteAddr().String()
	s.clients.Store(clientID, conn)
	defer s.clients.Delete(clientID)

	logger.Info("Client connected: %s", clientID)
	defer logger.Info("Client disconnected: %s", clientID)

	buffer := make([]byte, 4096)
	for {
		select {
		case <-s.quit:
			return
		default:
			n, err := conn.Read(buffer)
			if err != nil {
				logger.Error("Error reading from connection %s: %v", clientID, err)
				return
			}

			logger.Debug("Received %d bytes from %s", n, clientID)

			// Create a task for the worker pool
			task := Task{
				Data:    buffer[:n],
				ConnID:  clientID,
				Context: context.Background(),
			}

			s.workerPool.Submit(task)
		}
	}
}

// handleTask processes individual tasks
func (s *TCPServer) handleTask(task Task) error {
	logger.Debug("Processing task from client %s", task.ConnID)

	// Process the received data
	// This is where you would implement your specific protocol handling
	// For example, parsing the message and generating a response

	// Get the client connection
	if clientConn, ok := s.clients.Load(task.ConnID); ok {
		conn := clientConn.(net.Conn)

		// Echo the data back (replace with actual protocol handling)
		_, err := conn.Write(task.Data)
		if err != nil {
			logger.Error("Failed to write response to client %s: %v", task.ConnID, err)
			return fmt.Errorf("failed to write response: %v", err)
		}
		logger.Debug("Sent response to client %s", task.ConnID)
	} else {
		logger.Warn("Client %s not found", task.ConnID)
	}

	return nil
}

// Stop gracefully shuts down the server
func (s *TCPServer) Stop() {
	logger.Info("Stopping TCP server")
	close(s.quit)
	s.listener.Close()
	s.workerPool.Stop()

	// Close all client connections
	s.clients.Range(func(key, value interface{}) bool {
		conn := value.(net.Conn)
		clientID := key.(string)
		logger.Debug("Closing connection to client %s", clientID)
		conn.Close()
		return true
	})

	logger.Info("TCP server stopped")
}
