package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hello-world/logic-server/config"
	"github.com/hello-world/logic-server/entity"
	"github.com/hello-world/logic-server/mcp"
	"github.com/hello-world/logic-server/storage"
)

// Server 游戏服务器
type Server struct {
	storage  storage.Storage
	mcp      mcp.MCPClient
	mu       sync.RWMutex
	entities map[string]entity.Entity
	offline  bool // 是否为离线模式
}

// NewServer 创建服务器实例
func NewServer() (*Server, error) {
	var store storage.Storage
	var mcpClient mcp.MCPClient
	var err error

	// 初始化存储
	if config.GlobalConfig.Redis.Enabled {
		store, err = storage.NewRedisStorage(
			config.GlobalConfig.Redis.Host,
			config.GlobalConfig.Redis.Port,
			config.GlobalConfig.Redis.Password,
			config.GlobalConfig.Redis.DB,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Redis: %v\nPlease make sure Redis is running on %s:%d",
				err, config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port)
		}
	} else {
		log.Println("Running in offline mode: Redis storage is disabled")
		store = storage.NewMemoryStorage() // 使用内存存储
	}

	// 初始化MCP客户端
	if config.GlobalConfig.MCP.Enabled {
		mcpClient = mcp.NewMCPClient(
			config.GlobalConfig.MCP.Host,
			config.GlobalConfig.MCP.Port,
		)
	} else {
		log.Println("Running in offline mode: MCP client is disabled")
		mcpClient = mcp.NewMockMCPClient() // 使用模拟MCP客户端
	}

	return &Server{
		storage:  store,
		mcp:      mcpClient,
		entities: make(map[string]entity.Entity),
		offline:  !config.GlobalConfig.Redis.Enabled || !config.GlobalConfig.MCP.Enabled,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	if config.GlobalConfig.MCP.Enabled {
		// 启动MCP客户端
		if err := s.mcp.Start(); err != nil {
			return fmt.Errorf("failed to connect to MCP server: %v\nPlease make sure MCP server is running on %s:%d",
				err, config.GlobalConfig.MCP.Host, config.GlobalConfig.MCP.Port)
		}
		log.Printf("Successfully connected to MCP server at %s:%d", config.GlobalConfig.MCP.Host, config.GlobalConfig.MCP.Port)
	}

	if config.GlobalConfig.Redis.Enabled {
		log.Printf("Successfully connected to Redis at %s:%d", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port)
	}

	if s.offline {
		log.Printf("Server is running in offline mode")
	} else {
		log.Printf("Server is running in online mode")
	}

	// 启动主循环
	go s.mainLoop()

	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	if config.GlobalConfig.MCP.Enabled {
		return s.mcp.Close()
	}
	return nil
}

// mainLoop 主循环
func (s *Server) mainLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.update()
		}
	}
}

// update 更新逻辑
func (s *Server) update() {
	s.mu.RLock()
	entities := make([]entity.Entity, 0, len(s.entities))
	for _, e := range s.entities {
		entities = append(entities, e)
	}
	s.mu.RUnlock()

	// 向AI发送请求
	ctx := context.Background()
	response, err := s.mcp.RequestAI(ctx, entities)
	if err != nil {
		log.Printf("Error requesting AI: %v", err)
		return
	}

	// 处理AI响应
	s.processAIResponse(response)
}

// processAIResponse 处理AI响应
func (s *Server) processAIResponse(response map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 更新实体状态
	for id, data := range response {
		if e, ok := s.entities[id]; ok {
			if updates, ok := data.(map[string]interface{}); ok {
				for k, v := range updates {
					e.SetData(k, v)
				}
				// 保存到存储
				if err := s.storage.SaveEntity(context.Background(), e); err != nil {
					log.Printf("Error saving entity %s: %v", id, err)
				}
			}
		}
	}
}

// AddEntity 添加实体
func (s *Server) AddEntity(e entity.Entity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entities[e.GetID()] = e
}

// RemoveEntity 移除实体
func (s *Server) RemoveEntity(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entities, id)
}

// GetEntity 获取实体
func (s *Server) GetEntity(id string) (entity.Entity, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entities[id]
	return e, ok
}
