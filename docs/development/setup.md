# 开发环境设置指南

## 1. 环境要求
### 1.1 基础环境
- Go 1.20+
- Godot 4.4
- Redis 7.0+
- Git

### 1.2 开发工具
- Visual Studio Code (推荐)
- GoLand (可选)
- Godot Editor

## 2. 项目结构
```
project/
├── ai-mcp-server/     # AI MCP服务器
├── server/            # 业务服务器
├── game-client/       # 游戏客户端
└── docs/             # 项目文档
```

## 3. 安装步骤

### 3.1 AI MCP Server
```bash
# 克隆项目
git clone [项目地址]

# 进入项目目录
cd ai-mcp-server

# 安装依赖
go mod download

# 配置环境变量
export GEMINI_API_KEY=your_api_key

# 运行服务器
go run main.go
```

# 修改完ai api key后正常运行
DAP server listening at: 127.0.0.1:58627
Type 'dlv help' for list of commands.
Prompt: Hello! Please introduce yourself.
Response: Hello! I am a large language model, trained by Google. I am designed to provide information and complete tasks based on the instructions and data I've been trained on. Essentially, I can generate text, translate languages, write different kinds of creative content, and answer your questions in an informative way. I am still under development, learning new things every day!

How can I help you today?

Process 37492 has exited with status 0
Detaching

### 3.2 业务服务器
```bash
# 进入项目目录
cd server

# 安装依赖
go mod download

# 运行服务器
go run main.go
```

### 3.3 游戏客户端
1. 安装Godot 4.4
2. 打开Godot编辑器
3. 导入项目
4. 打开game-client目录
5. 运行项目

### 3.4 Redis设置
```bash
# 安装Redis
# Windows: 下载Redis安装包并安装
# Linux: sudo apt-get install redis-server

# 配置Redis
# 修改redis.conf
port 6666

# 启动Redis
redis-server redis.conf
```

## 4. 开发工作流

### 4.1 代码规范
- Go代码遵循标准Go代码规范
- C#代码遵循C#编码规范
- 使用gofmt进行Go代码格式化
- 使用EditorConfig统一编辑器配置

### 4.2 版本控制
- 使用Git进行版本控制
- 遵循Git Flow工作流
- 提交信息遵循约定式提交规范

### 4.3 测试
- Go测试使用标准testing包
- 单元测试覆盖率要求>80%
- 集成测试使用测试环境

## 5. 调试指南

### 5.1 AI MCP Server调试
- 使用Delve进行Go调试
- 设置日志级别为DEBUG
- 使用pprof进行性能分析

### 5.2 业务服务器调试
- 使用Delve进行Go调试
- 启用详细日志
- 使用pprof进行性能分析

### 5.3 游戏客户端调试
- 使用Godot调试器
- 启用开发者控制台
- 使用性能监视器

## 6. 常见问题

### 6.1 Redis连接问题
- 检查Redis服务是否运行
- 验证端口配置
- 检查防火墙设置

### 6.2 API调用问题
- 验证API密钥
- 检查网络连接
- 查看错误日志

### 6.3 游戏客户端问题
- 检查Godot版本兼容性
- 验证资源文件完整性
- 检查网络连接状态 