package ai

import (
	"context"
	"log"
	"time"

	"github.com/hello-world/common/proto/ai"
	"github.com/hello-world/logic-server/entity"
)

// Client represents an AI service client
type Client struct {
	grpcPool *grpcpool.GrpcPool
}

// NewClient creates a new AI service client
func NewClient(target string, poolSize int) *Client {
	return &Client{
		grpcPool: grpcpool.NewPool(target, poolSize),
	}
}

// GetDecision requests an AI decision for a woodcutter
func (c *Client) GetDecision(ctx context.Context, woodcutter *entity.Woodcutter, entities []interface{}) (*ai.DecisionResponse, error) {
	// Get connection from grpcPool
	conn, err := c.grpcPool.Get()
	if err != nil {
		return nil, err
	}
	defer c.grpcPool.Put(conn)

	// Create AI service client
	client := ai.NewAIServiceClient(conn)

	// Convert woodcutter data to proto format
	woodcutterData := woodcutter.GetData()
	actorInfo := &ai.EntityInfo{
		Id:   woodcutterData["id"].(string),
		Type: woodcutterData["type"].(string),
		Position: &ai.Position{
			X: float32(woodcutterData["position"].(entity.Position).X),
			Y: float32(woodcutterData["position"].(entity.Position).Y),
			Z: float32(woodcutterData["position"].(entity.Position).Z),
		},
		Attributes: map[string]float32{
			"hp":     float32(woodcutterData["hp"].(float64)),
			"energy": float32(woodcutterData["energy"].(float64)),
		},
	}

	// Convert entities to proto format
	entityInfos := make([]*ai.EntityInfo, 0, len(entities))
	for _, e := range entities {
		if data, ok := e.(map[string]interface{}); ok {
			pos := data["position"].(entity.Position)
			entityInfo := &ai.EntityInfo{
				Id:   data["id"].(string),
				Type: data["type"].(string),
				Position: &ai.Position{
					X: float32(pos.X),
					Y: float32(pos.Y),
					Z: float32(pos.Z),
				},
			}
			if hp, ok := data["hp"].(float64); ok {
				entityInfo.Attributes = map[string]float32{"hp": float32(hp)}
			}
			entityInfos = append(entityInfos, entityInfo)
		}
	}

	// Create request
	request := &ai.DecisionRequest{
		Actor:        actorInfo,
		MapInfo:      &ai.MapInfo{Entities: entityInfos},
		CurrentState: string(woodcutter.State),
	}

	// Get decision from AI service
	// Create a timeout context for the RPC call
	rpcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := client.GetDecision(rpcCtx, request)
	if err != nil {
		if rpcCtx.Err() == context.DeadlineExceeded {
			log.Printf("AI service timeout: %v", err)
		} else {
			log.Printf("Failed to get AI decision: %v", err)
		}
		return nil, err
	}

	return response, nil
}

// Close closes the client's connection grpcPool
func (c *Client) Close() {
	c.grpcPool.Close()
}
