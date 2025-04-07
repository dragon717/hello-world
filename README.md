# 游戏项目文档

## 项目概述
这是一个基于分布式架构的游戏项目,包含AI MCP服务器、业务服务器和游戏客户端。

## 技术架构
### 服务器端
- **AI MCP Server**
  - 语言: Go
  - 功能: 负责AI相关功能,集成Gemini API
  - 职责: 处理AI决策、对话、行为等

- **业务服务器 (Server)**
  - 语言: Go
  - 功能: 处理游戏核心业务逻辑
  - 职责: 玩家管理、游戏状态同步、业务规则处理

### 客户端
- **游戏客户端 (Gameplay)**
  - 引擎: Godot
  - 语言: C#
  - 功能: 游戏表现层、玩家交互

### 数据存储
- **数据库**
  - 类型: Redis
  - 端口: 6666
  - 存储格式: Key-Value (Value为JSON格式)

### 通信协议
- 协议格式: JSON
- 通信方式: TCP

## 文档导航
- [系统架构](./docs/architecture/system-design.md)
- [游戏设计](./docs/game-design/game-overview.md)
- [开发指南](./docs/development/setup.md)

## 快速开始
1. 环境要求
   - Go 1.20+
   - Godot 4.4
   - Redis 7.0+
   - Gemini API Key

2. 安装步骤
   - [AI MCP Server 安装指南](./docs/development/setup.md#ai-mcp-server)
   - [业务服务器安装指南](./docs/development/setup.md#business-server)
   - [游戏客户端安装指南](./docs/development/setup.md#game-client)

3. 运行方式
   - 启动Redis服务 (端口6666)
   - 启动AI MCP Server
   - 启动业务服务器
   - 运行游戏客户端

## 开发状态
当前版本：v0.0.1 (开发中)

## 许可证
MIT License