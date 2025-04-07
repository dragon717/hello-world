package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hello-world/logic-server/entity"
)

// MCPClient MCP客户端接口
type MCPClient interface {
	Start() error
	Close() error
	RequestAI(ctx context.Context, entities []entity.Entity) (map[string]interface{}, error)
}

// WebSocketMCPClient WebSocket MCP客户端实现
type WebSocketMCPClient struct {
	conn *websocket.Conn
	host string
	port int
	done chan struct{}
}

// NewMCPClient 创建MCP客户端
func NewMCPClient(host string, port int) MCPClient {
	return &WebSocketMCPClient{
		host: host,
		port: port,
		done: make(chan struct{}),
	}
}

// Connect 连接到MCP服务器
func (c *WebSocketMCPClient) Connect() error {
	url := fmt.Sprintf("ws://%s:%d/ws", c.host, c.port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Close 关闭连接
func (c *WebSocketMCPClient) Close() error {
	close(c.done)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// RequestAI 向AI发送请求
func (c *WebSocketMCPClient) RequestAI(ctx context.Context, entities []entity.Entity) (map[string]interface{}, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected to MCP server")
	}

	// 准备请求数据
	request := map[string]interface{}{
		"entities":  entities,
		"timestamp": time.Now().Unix(),
	}

	// 发送请求
	if err := c.conn.WriteJSON(request); err != nil {
		return nil, err
	}

	// 等待响应
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// Start 启动客户端
func (c *WebSocketMCPClient) Start() error {
	if err := c.Connect(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				// 保持连接活跃
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Error sending ping: %v", err)
					return
				}
				time.Sleep(30 * time.Second)
			}
		}
	}()

	return nil
}
