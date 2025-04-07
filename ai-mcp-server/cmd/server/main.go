package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bugod/ai-mcp-server/internal/config"
	"github.com/bugod/ai-mcp-server/internal/service"
	"github.com/gorilla/websocket"
)

// Request represents the incoming request structure
type Request struct {
	Entities  []interface{} `json:"entities"`
	Timestamp int64         `json:"timestamp"`
}

// Response represents the outgoing response structure
type Response struct {
	Success bool                  `json:"success"`
	Data    map[string]ActionData `json:"data"`
	Error   string                `json:"error,omitempty"`
}

// ActionData represents the AI decision for an entity
type ActionData struct {
	ActionType string   `json:"action_type"`
	TargetID   string   `json:"target_id,omitempty"`
	Position   Position `json:"position,omitempty"`
	State      string   `json:"state,omitempty"`
}

// Position represents a 3D position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

func main() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Load configuration from the same directory
	configPath := filepath.Join(wd, "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Gemini service
	geminiService, err := service.NewGeminiService(cfg.Gemini.APIKey, cfg.Gemini.Model)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}
	defer geminiService.Close()

	// 设置 WebSocket 路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, geminiService)
	})

	// 启动 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("AI MCP Server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, geminiService *service.GeminiService) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("New WebSocket connection from %s", clientAddr)

	// 设置读写超时
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

	// 处理消息循环
	for {
		// 读取消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from %s: %v", clientAddr, err)
			return
		}

		// 重置读写超时
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

		// 解析请求
		var request Request
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Error parsing request: %v", err)
			sendErrorResponse(conn, "Invalid request format")
			continue
		}

		// 为每个实体生成决策
		response := Response{
			Success: true,
			Data:    make(map[string]ActionData),
		}

		// 处理每个实体
		for _, entity := range request.Entities {
			entityData, ok := entity.(map[string]interface{})
			if !ok {
				continue
			}

			id, _ := entityData["id"].(string)
			entityType, _ := entityData["type"].(string)

			// 根据实体类型生成决策
			switch entityType {
			case "WOODCUTTER":
				// 生成樵夫的决策
				ctx := context.Background()
				prompt := fmt.Sprintf("As an AI, generate a decision for a woodcutter. Current state: %v", entityData)
				aiResponse, err := geminiService.GenerateResponse(ctx, prompt)
				if err != nil {
					log.Printf("Error generating response: %v", err)
					continue
				}

				// 解析 AI 响应并生成决策
				decision := parseAIResponse(aiResponse)
				response.Data[id] = decision
			}
		}

		// 发送响应
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Error writing to %s: %v", clientAddr, err)
			return
		}
	}
}

func sendErrorResponse(conn *websocket.Conn, message string) {
	response := Response{
		Success: false,
		Error:   message,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}

func parseAIResponse(aiResponse string) ActionData {
	// 这里可以添加更复杂的解析逻辑
	// 现在先返回一个简单的决策
	return ActionData{
		ActionType: "SEARCH",
		State:      "SEARCHING",
	}
}
