package mcp

import (
	"context"
	"math/rand"
	"time"

	"github.com/hello-world/logic-server/entity"
)

// MockMCPClient 模拟MCP客户端
type MockMCPClient struct {
	done chan struct{}
}

// NewMockMCPClient 创建模拟MCP客户端
func NewMockMCPClient() *MockMCPClient {
	return &MockMCPClient{
		done: make(chan struct{}),
	}
}

// Start 启动模拟客户端
func (c *MockMCPClient) Start() error {
	return nil
}

// Close 关闭模拟客户端
func (c *MockMCPClient) Close() error {
	close(c.done)
	return nil
}

// RequestAI 模拟AI请求
func (c *MockMCPClient) RequestAI(ctx context.Context, entities []entity.Entity) (map[string]interface{}, error) {
	// 模拟处理延迟
	time.Sleep(10 * time.Millisecond)

	response := make(map[string]interface{})

	// 为每个实体生成随机更新
	for _, e := range entities {
		updates := make(map[string]interface{})

		// 根据实体类型生成不同的模拟数据
		switch e.GetType() {
		case entity.EntityTypeSystem:
			updates["timestamp"] = time.Now().Unix()
			updates["status"] = "running"

		case entity.EntityTypeUser:
			updates["position"] = map[string]float64{
				"x": rand.Float64() * 100,
				"y": rand.Float64() * 100,
			}
			updates["health"] = 80 + rand.Float64()*20

		case entity.EntityTypeNPC:
			updates["position"] = map[string]float64{
				"x": rand.Float64() * 100,
				"y": rand.Float64() * 100,
			}
			updates["state"] = []string{"idle", "walking", "running"}[rand.Intn(3)]
		}

		response[e.GetID()] = updates
	}

	return response, nil
}
