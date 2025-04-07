# Logic Server

Logic Server是一个游戏业务服务器，负责管理游戏实体状态并与AI MCP服务器通信。

## 功能特点

- 实体管理：支持系统级、用户级和NPC实体的管理
- 数据存储：使用Redis作为持久化存储
- AI集成：与AI MCP服务器集成，实现智能决策
- 实时更新：支持实体的实时状态更新

## 项目结构

```
logic-server/
├── config/         # 配置文件
├── entity/         # 实体定义
├── mcp/           # MCP客户端
├── server/        # 服务器逻辑
├── storage/       # 存储接口
└── main.go        # 程序入口
```

## 实体类型

- SystemEntity：系统级实体
- UserEntity：用户级实体
- NPCEntity：NPC实体

## 配置说明

配置文件位于 `config/config.go`，包含以下配置项：

- Server：服务器配置
  - Host：监听地址
  - Port：监听端口
- Redis：Redis配置
  - Host：Redis地址
  - Port：Redis端口
  - Password：Redis密码
  - DB：数据库编号
- MCP：MCP服务器配置
  - Host：MCP服务器地址
  - Port：MCP服务器端口

## 运行要求

- Go 1.21+
- Redis 6.0+
- AI MCP服务器

## 运行方法

1. 确保Redis和MCP服务器已启动
2. 运行服务器：
   ```bash
   go run main.go
   ```

## 开发说明

1. 实体数据以JSON格式存储
2. 实体状态更新通过AI MCP服务器进行
3. 所有实体操作都是线程安全的 