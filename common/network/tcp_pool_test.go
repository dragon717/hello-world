package network

import (
	"net"
	"testing"
	"time"
)

// mockListener implements a mock TCP listener for testing
type mockListener struct {
	addr net.Addr
	ch   chan net.Conn
}

func newMockListener(addr net.Addr) *mockListener {
	return &mockListener{
		addr: addr,
		ch:   make(chan net.Conn),
	}
}

func (l *mockListener) Accept() (net.Conn, error) {
	conn := <-l.ch
	return conn, nil
}

func (l *mockListener) Close() error {
	close(l.ch)
	return nil
}

func (l *mockListener) Addr() net.Addr {
	return l.addr
}

// mockConn implements a mock TCP connection for testing
type mockConn struct {
	net.Conn
	isClosed bool
}

func newMockConn() *mockConn {
	return &mockConn{}
}

func (c *mockConn) Close() error {
	c.isClosed = true
	return nil
}

func TestNewTCPPool(t *testing.T) {
	tests := []struct {
		name    string
		address string
		maxSize int
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "Valid configuration",
			address: "localhost:8080",
			maxSize: 10,
			timeout: 30 * time.Second,
			wantErr: false,
		},
		{
			name:    "Invalid address",
			address: "invalid:address",
			maxSize: 10,
			timeout: 30 * time.Second,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := NewTCPPool(tt.address, tt.maxSize, tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTCPPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pool == nil {
				t.Error("NewTCPPool() returned nil pool")
			}
		})
	}
}

func TestTCPPool_GetPut(t *testing.T) {
	pool, err := NewTCPPool("localhost:8080", 2, 30*time.Second)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}

	// Test Get
	conn1, err := pool.Get()
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if conn1 == nil {
		t.Error("Get() returned nil connection")
	}

	// Test Put
	pool.Put(conn1)

	// Test multiple Get/Put operations
	conn2, _ := pool.Get()
	conn3, _ := pool.Get()

	// GrpcPool is now empty
	pool.Put(conn2)
	pool.Put(conn3)
}

func TestTCPPool_Monitor(t *testing.T) {
	pool, err := NewTCPPool("localhost:8080", 2, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}

	// Get a connection and let it expire
	conn, err := pool.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	pool.Put(conn)

	// Wait for the connection to expire
	time.Sleep(2 * time.Second)

	// Try to get the connection again
	conn2, err := pool.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if conn2 == conn {
		t.Error("Got same expired connection")
	}
}

func TestTCPPool_Close(t *testing.T) {
	// Create a mock listener and connection
	mockConn1 := newMockConn()
	mockConn2 := newMockConn()

	pool := &TCPPool{
		connections: make(chan *TCPConn, 2),
		maxSize:     2,
		address:     "localhost:8080",
		timeout:     30 * time.Second,
	}

	// Add connections to the pool
	pool.connections <- &TCPConn{
		conn:       mockConn1,
		createTime: time.Now(),
		lastUsed:   time.Now(),
	}
	pool.connections <- &TCPConn{
		conn:       mockConn2,
		createTime: time.Now(),
		lastUsed:   time.Now(),
	}

	// Close the pool
	pool.Close()

	// Verify that the connections were closed
	if !mockConn1.isClosed {
		t.Error("First connection was not closed")
	}
	if !mockConn2.isClosed {
		t.Error("Second connection was not closed")
	}

	// Verify that getting a new connection fails
	_, err := pool.Get()
	if err == nil {
		t.Error("Expected error after closing pool")
	} else if err.Error() != "pool is closed" {
		t.Errorf("Expected 'pool is closed' error, got: %v", err)
	}
}
