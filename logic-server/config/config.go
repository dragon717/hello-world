package config

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
		Enabled  bool   `json:"enabled"`
	} `json:"redis"`

	MCP struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		Enabled bool   `json:"enabled"`
	} `json:"mcp"`
}

var GlobalConfig Config

func init() {
	// 默认配置
	GlobalConfig.Server.Host = "0.0.0.0"
	GlobalConfig.Server.Port = 8080

	GlobalConfig.Redis.Host = "localhost"
	GlobalConfig.Redis.Port = 6379
	GlobalConfig.Redis.DB = 0
	GlobalConfig.Redis.Enabled = false

	GlobalConfig.MCP.Host = "localhost"
	GlobalConfig.MCP.Port = 8082
	GlobalConfig.MCP.Enabled = false
}
