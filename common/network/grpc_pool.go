package network

import (
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// GrpcPool represents a pool of gRPC connections
type GrpcPool struct {
	mu          sync.Mutex
	connections []*grpc.ClientConn
	target      string
	maxSize     int
	currentSize int
}

// NewGrpcPool creates a new connection pool
func NewGrpcPool(target string, maxSize int) *GrpcPool {
	return &GrpcPool{
		connections: make([]*grpc.ClientConn, 0, maxSize),
		target:      target,
		maxSize:     maxSize,
	}
}

// Get retrieves a connection from the pool
func (p *GrpcPool) Get() (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If there are available connections, return one
	if len(p.connections) > 0 {
		conn := p.connections[len(p.connections)-1]
		p.connections = p.connections[:len(p.connections)-1]
		return conn, nil
	}

	// If we can create a new connection, do so
	if p.currentSize < p.maxSize {
		conn, err := grpc.Dial(
			p.target,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(5*time.Second),
		)
		if err != nil {
			return nil, err
		}
		p.currentSize++
		return conn, nil
	}

	return nil, errors.New("connection pool exhausted")
}

// Put returns a connection to the pool
func (p *GrpcPool) Put(conn *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.connections = append(p.connections, conn)
}

// Close closes all connections in the pool
func (p *GrpcPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.connections {
		conn.Close()
	}
	p.connections = nil
	p.currentSize = 0
}
