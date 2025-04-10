package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hello-world/logic-server/entity"
	"github.com/hello-world/logic-server/server"
)

func main() {
	// 创建服务器实例
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	// 创建测试实体
	system := entity.NewSystemEntity("sys1")
	user := entity.NewUserEntity("user1", "TestUser")
	npc := entity.NewNPCEntity("npc1", "TestNPC")

	// 添加实体
	srv.AddEntity(system)
	srv.AddEntity(user)
	srv.AddEntity(npc)

	// 定期打印实体状态
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if e, ok := srv.GetEntity("user1"); ok {
					data := e.GetData()
					log.Printf("User1 state: %+v", data)
				}
				if e, ok := srv.GetEntity("npc1"); ok {
					data := e.GetData()
					log.Printf("NPC1 state: %+v", data)
				}
			}
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
}
