package network

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// TCPConn wraps a net.Conn with additional metadata
type TCPConn struct {
	conn       net.Conn
	lastUsed   time.Time
	createTime time.Time
	inUse      bool
}

// TCPPool represents a pool of TCP connections
type TCPPool struct {
	mu          sync.RWMutex
	connections chan *TCPConn
	maxSize     int
	address     string
	timeout     time.Duration
	closed      bool
}

// NewTCPPool creates a new TCP connection pool
func NewTCPPool(address string, maxSize int, timeout time.Duration) (*TCPPool, error) {
	// Try to establish a test connection to validate the address
	testConn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to validate address: %v", err)
	}
	testConn.Close()

	pool := &TCPPool{
		connections: make(chan *TCPConn, maxSize),
		maxSize:     maxSize,
		address:     address,
		timeout:     timeout,
	}

	// Start connection monitor
	go pool.monitor()

	return pool, nil
}

// Get retrieves a connection from the pool
func (p *TCPPool) Get() (*TCPConn, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("pool is closed")
	}
	p.mu.RUnlock()

	select {
	case conn, ok := <-p.connections:
		if !ok {
			return nil, fmt.Errorf("pool is closed")
		}
		if time.Since(conn.lastUsed) > p.timeout {
			conn.conn.Close()
			return p.createConnection()
		}
		conn.inUse = true
		conn.lastUsed = time.Now()
		return conn, nil
	default:
		return p.createConnection()
	}
}

// Put returns a connection to the pool
func (p *TCPPool) Put(conn *TCPConn) {
	p.mu.RLock()
	if p.closed {
		conn.conn.Close()
		p.mu.RUnlock()
		return
	}
	p.mu.RUnlock()

	conn.inUse = false
	conn.lastUsed = time.Now()

	// If the pool is full, close the connection
	select {
	case p.connections <- conn:
	default:
		conn.conn.Close()
	}
}

// createConnection establishes a new TCP connection
func (p *TCPPool) createConnection() (*TCPConn, error) {
	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		return nil, err
	}

	return &TCPConn{
		conn:       conn,
		createTime: time.Now(),
		lastUsed:   time.Now(),
	}, nil
}

// monitor periodically checks and removes stale connections
func (p *TCPPool) monitor() {
	ticker := time.NewTicker(p.timeout / 2)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.RLock()
		if p.closed {
			p.mu.RUnlock()
			return
		}
		p.mu.RUnlock()

		currentConns := len(p.connections)
		for i := 0; i < currentConns; i++ {
			select {
			case conn := <-p.connections:
				if time.Since(conn.lastUsed) > p.timeout {
					conn.conn.Close()
				} else {
					p.connections <- conn
				}
			default:
				break
			}
		}
	}
}

// Close closes all connections in the pool
func (p *TCPPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	// Close all connections in the pool
	close(p.connections)
	for conn := range p.connections {
		if conn != nil && conn.conn != nil {
			conn.conn.Close()
		}
	}
}
