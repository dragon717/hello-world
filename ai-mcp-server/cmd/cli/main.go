package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugod/ai-mcp-server/internal/config"
	"github.com/bugod/ai-mcp-server/internal/service"
)

func main() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	// 解析命令行参数
	configPath := flag.String("config", filepath.Join(wd, "config.yaml"), "Path to config file")
	flag.Parse()
	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// 创建Gemini服务
	geminiService, err := service.NewGeminiService(cfg.Gemini.APIKey, cfg.Gemini.Model)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}
	defer geminiService.Close()
	// 创建上下文
	ctx := context.Background()
	// 创建读取器
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("AI MCP Server CLI")
	fmt.Println("Type 'exit' to quit")
	fmt.Println("-------------------")
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		// 去除换行符
		input = strings.TrimSpace(input)
		// 检查退出命令
		if input == "exit" {
			break
		}
		// 跳过空输入
		if input == "" {
			continue
		}
		// 生成响应
		response, err := geminiService.GenerateResponse(ctx, input)
		if err != nil {
			log.Printf("Error generating response: %v", err)
			continue
		}
		fmt.Println("Response:", response)
		fmt.Println("-------------------")
	}
}
