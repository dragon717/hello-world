package service

import (
	"context"
	"math"
	"strings"

	ai "github.com/hello-world/common/protocol"
	"google.golang.org/grpc"
)

// AIService implements the AI service
type AIService struct {
	ai.UnimplementedAIServiceServer
	geminiService *GeminiService
}

// NewAIService creates a new AI service
func NewAIService(geminiService *GeminiService) *AIService {
	return &AIService{
		geminiService: geminiService,
	}
}

// RegisterService registers the AI service with a gRPC server
func (s *AIService) RegisterService(server *grpc.Server) {
	ai.RegisterAIServiceServer(server, s)
}

// GetDecision implements the GetDecision RPC method
func (s *AIService) GetDecision(ctx context.Context, req *ai.DecisionRequest) (*ai.DecisionResponse, error) {
	// Based on the current state, generate appropriate prompts for Gemini
	var prompt string
	switch req.CurrentState {
	case "IDLE":
		prompt = "As a woodcutter, what should be my first action: learning woodcutting skill or searching for trees?"
	case "LEARNING":
		prompt = "I am learning woodcutting skills. What should be my next action?"
	case "SEARCHING":
		// Find nearest tree
		var nearestTree *ai.EntityInfo
		var minDist float32 = math.MaxFloat32
		for _, entity := range req.MapInfo.Entities {
			if entity.Type == "TREE" {
				dist := distance(req.Actor.Position, entity.Position)
				if dist < minDist {
					minDist = dist
					nearestTree = entity
				}
			}
		}
		if nearestTree != nil {
			prompt = "I found a tree. Should I move to it and start chopping?"
		} else {
			prompt = "I can't find any trees nearby. What should I do?"
		}
	case "MOVING":
		prompt = "I am moving towards a tree. Should I continue moving or start chopping?"
	case "CHOPPING":
		prompt = "I am chopping a tree. Should I continue chopping or do something else?"
	}

	// Get AI response from Gemini
	response, err := s.geminiService.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Parse Gemini response and convert to appropriate action
	// This is a simplified version - in reality, you'd want more sophisticated parsing
	var actionType string
	var targetId string
	var targetPosition *ai.Position
	var subCommands []string

	// Simple parsing based on keywords in response
	switch {
	case contains(response, "learn", "skill"):
		actionType = "LEARN_SKILL"
		subCommands = []string{"learn_woodcutting"}
	case contains(response, "search", "tree"):
		actionType = "SEARCH"
		subCommands = []string{"scan_surroundings"}

	default:
		actionType = "IDLE"
		subCommands = []string{"wait"}
	}

	return &ai.DecisionResponse{
		ActionType:     actionType,
		TargetId:       targetId,
		TargetPosition: targetPosition,
		SubCommands:    subCommands,
	}, nil
}

// Helper functions

func distance(p1, p2 *ai.Position) float32 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(strings.ToLower(s), strings.ToLower(substr)) {
			return true
		}
	}
	return false
}
