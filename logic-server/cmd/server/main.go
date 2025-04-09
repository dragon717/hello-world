package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hello-world/logic-server/entity"
	"github.com/hello-world/logic-server/offline_ai"
)

const (
	aiServerAddress = "localhost:8888"
	aiPoolSize      = 5
)

func main() {
	log.Println("Starting woodcutter simulation...")

	// 创建 AI 客户端
	aiClient := offline_ai.NewClient(aiServerAddress, aiPoolSize)
	defer aiClient.Close()

	// 创建地图和实体
	woodcutter := entity.NewWoodcutter("wc1", entity.Position{X: 0, Y: 0, Z: 0})
	trees := []interface{}{
		entity.NewTree("tree1", entity.Position{X: 10, Y: 0, Z: 0}),
		entity.NewTree("tree2", entity.Position{X: 20, Y: 5, Z: 0}),
		entity.NewTree("tree3", entity.Position{X: -15, Y: -5, Z: 0}),
	}

	// 获取所有实体（用于发送给 AI 服务）
	allEntities := []interface{}{woodcutter}
	for _, t := range trees {
		allEntities = append(allEntities, t)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动樵夫逻辑处理循环
	go processWoodcutterLogic(ctx, aiClient, woodcutter, trees)

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

// processWoodcutterLogic 处理樵夫的逻辑
func processWoodcutterLogic(ctx context.Context, aiClient *offline_ai.Client, woodcutter *entity.Woodcutter, trees []interface{}) {
	// 循环处理樵夫生命周期
	ticker := time.NewTicker(2 * time.Second) // 每2秒进行一次决策
	defer ticker.Stop()

	// 获取所有实体（用于发送给 AI 服务）
	allEntities := []interface{}{woodcutter}
	for _, t := range trees {
		allEntities = append(allEntities, t)
	}

	// 樵夫生命周期主循环
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 获取AI决策
			response, err := aiClient.GetDecision(ctx, woodcutter, allEntities)
			if err != nil {
				log.Printf("Failed to get AI decision: %v", err)
				continue
			}

			log.Printf("AI Decision: Action=%s, Target=%s", response.ActionType, response.TargetId)

			// 根据决策执行相应动作
			switch response.ActionType {
			case "LEARN_SKILL":
				woodcutter.UpdateState(entity.StateLearning)
				woodcutter.LearnSkill("woodcutting")
				log.Printf("Woodcutter learned woodcutting skill")
				woodcutter.UpdateState(entity.StateIdle)

			case "SEARCH":
				woodcutter.UpdateState(entity.StateSearching)
				log.Printf("Woodcutter is searching for trees")

			case "MOVE":
				woodcutter.UpdateState(entity.StateMoving)
				if response.TargetPosition != nil {
					woodcutter.MoveTo(entity.Position{
						X: float64(response.TargetPosition.X),
						Y: float64(response.TargetPosition.Y),
						Z: float64(response.TargetPosition.Z),
					})
					log.Printf("Woodcutter moved to position (%f, %f, %f)",
						woodcutter.Position.X, woodcutter.Position.Y, woodcutter.Position.Z)
				}

			case "CHOP":
				woodcutter.UpdateState(entity.StateChopping)
				// 检查是否有砍树技能
				if !woodcutter.HasSkill("woodcutting") {
					log.Printf("Woodcutter doesn't have woodcutting skill!")
					woodcutter.UpdateState(entity.StateIdle)
					continue
				}

				// 找到目标树
				var targetTree *entity.Tree
				for _, t := range trees {
					tree := t.(*entity.Tree)
					if tree.ID == response.TargetId {
						targetTree = tree
						break
					}
				}

				if targetTree == nil {
					log.Printf("Target tree not found: %s", response.TargetId)
					woodcutter.UpdateState(entity.StateIdle)
					continue
				}

				// 砍树
				log.Printf("Chopping tree %s", targetTree.ID)
				isDead := targetTree.TakeDamage(20) // 每次砍伤害20点
				log.Printf("Tree HP: %f", targetTree.HP)

				// 如果树死了
				if isDead {
					log.Printf("Tree %s has been chopped down!", targetTree.ID)
					// 获得木材
					woodcutter.AddToInventory("wood", targetTree.Wood)
					log.Printf("Woodcutter got %d wood, total: %d",
						targetTree.Wood, woodcutter.Inventory["wood"])

					// 从树列表中移除
					for i, t := range trees {
						tree := t.(*entity.Tree)
						if tree.ID == targetTree.ID {
							trees = append(trees[:i], trees[i+1:]...)
							break
						}
					}

					// 更新AI决策用的实体列表
					allEntities = []interface{}{woodcutter}
					for _, t := range trees {
						allEntities = append(allEntities, t)
					}

					woodcutter.UpdateState(entity.StateIdle)
				}

			default:
				woodcutter.UpdateState(entity.StateIdle)
				log.Printf("Woodcutter is idle")
			}
		}
	}
}
