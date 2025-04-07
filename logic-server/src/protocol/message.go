package protocol

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of message being sent
type MessageType int

const (
	// Message types
	TypeRequest MessageType = iota
	TypeResponse
	TypeError
)

// Message represents the basic message structure
type Message struct {
	Type      MessageType     `json:"type"`
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// Request represents a request message
type Request struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// Response represents a response message
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error,omitempty"`
}

// NewRequest creates a new request message
func NewRequest(action string, data interface{}) (*Message, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reqData, err := json.Marshal(Request{
		Action: action,
		Data:   dataBytes,
	})
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      TypeRequest,
		ID:        generateID(),
		Timestamp: time.Now(),
		Payload:   reqData,
	}, nil
}

// NewResponse creates a new response message
func NewResponse(requestID string, success bool, data interface{}, errMsg string) (*Message, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	respData, err := json.Marshal(Response{
		Success: success,
		Data:    dataBytes,
		Error:   errMsg,
	})
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      TypeResponse,
		ID:        requestID,
		Timestamp: time.Now(),
		Payload:   respData,
	}, nil
}

// generateID generates a unique message ID
func generateID() string {
	return time.Now().Format("20060102150405.000") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
