package config

import (
	"encoding/json" // JSON编解码包，用于读取和写入JSON格式的配置文件
	"log"           // 日志记录包，用于输出日志信息
	"os"            // 操作系统功能包，用于文件操作
)

// Config 应用配置 - 这是应用程序的完整配置结构
// 它包含了服务器、数据库和日志三个主要部分的配置
type Config struct {
	Server   ServerConfig   `json:"server"`   // 服务器相关配置
	Database DatabaseConfig `json:"database"` // 数据库相关配置
	Logging  LoggingConfig  `json:"logging"`  // 日志相关配置
}

// ServerConfig 服务器配置 - 定义Web服务器的运行参数
type ServerConfig struct {
	Port           string   `json:"port"`            // 服务器监听的端口号，如 "8080"
	Debug          bool     `json:"debug"`           // 是否启用调试模式，true时可能输出更多信息
	AllowedOrigins []string `json:"allowed_origins"` // CORS允许的来源，用于跨域请求控制
	RateLimit      int      `json:"rate_limit"`      // 速率限制，单位时间内允许的最大请求数
}

// DatabaseConfig 数据库配置 - 定义数据库连接参数
type DatabaseConfig struct {
	Type     string `json:"type"`     // 数据库类型，如 "mysql", "postgres", "memory"（内存数据库）
	Host     string `json:"host"`     // 数据库服务器主机名或IP地址
	Port     int    `json:"port"`     // 数据库服务器端口号
	Name     string `json:"name"`     // 数据库名称
	Username string `json:"username"` // 数据库用户名
	Password string `json:"password"` // 数据库密码
}

// LoggingConfig 日志配置 - 定义日志记录的行为和参数
type LoggingConfig struct {
	Level      string `json:"level"`       // 日志级别：debug, info, warn, error等
	File       string `json:"file"`        // 日志文件路径，如 "logs/app.log"
	MaxSize    int    `json:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `json:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `json:"max_age"`     // 日志文件保留的最大天数
}

// LoadConfig 加载配置
// 这个函数尝试从config.json文件加载配置，如果文件不存在或读取失败，则使用默认配置
// 工作流程：
// 1. 首先创建包含默认值的配置对象
// 2. 检查是否存在config.json文件
// 3. 如果存在，读取并解析该文件
// 4. 如果文件不存在或解析失败，使用默认配置
// 5. 返回配置对象
func LoadConfig() *Config {
	// 创建默认配置对象
	// 这是当没有配置文件或配置文件读取失败时使用的配置
	config := &Config{
		Server: ServerConfig{
			Port:           "8080",        // 默认监听8080端口
			Debug:          false,         // 默认关闭调试模式
			AllowedOrigins: []string{"*"}, // 默认允许所有来源（开发环境方便，生产环境应限制）
			RateLimit:      100,           // 默认每秒100个请求的速率限制
		},
		Database: DatabaseConfig{
			Type:     "memory",      // 默认使用内存数据库（无需安装外部数据库）
			Host:     "localhost",   // 默认数据库主机
			Port:     0,             // 默认端口0（通常表示使用默认端口或不需要端口）
			Name:     "xstreamtool", // 默认数据库名称
			Username: "",            // 默认无用户名
			Password: "",            // 默认无密码
		},
		Logging: LoggingConfig{
			Level:      "info",         // 默认日志级别：info（记录info及以上级别）
			File:       "logs/app.log", // 默认日志文件路径
			MaxSize:    10,             // 默认每个日志文件最大10MB
			MaxBackups: 5,              // 默认保留5个旧日志文件
			MaxAge:     30,             // 默认日志文件保留30天
		},
	}

	// 尝试从配置文件加载
	// 首先检查配置文件是否存在
	// os.Stat返回文件信息，如果文件不存在则返回错误
	if _, err := os.Stat("config.json"); err == nil {
		// 文件存在，读取文件内容
		data, err := os.ReadFile("config.json")
		if err != nil {
			// 读取文件失败，记录警告但继续使用默认配置
			// 这是"优雅降级"的设计：即使配置读取失败，应用也能启动
			log.Printf("⚠️ 读取配置文件失败: %v", err)
			return config
		}

		// 解析JSON文件内容到config结构体
		// 注意：这里使用了json.Unmarshal将JSON数据填充到已有的config对象中
		// JSON中的字段会覆盖默认值，JSON中没有的字段保持默认值
		if err := json.Unmarshal(data, config); err != nil {
			// JSON解析失败，记录警告但继续使用默认配置
			log.Printf("⚠️ 解析配置文件失败: %v", err)
		}
	}
	// 注意：如果config.json文件不存在，不会记录错误，直接使用默认配置
	// 这是有意为之的，让应用在首次运行时能自动使用默认配置启动

	return config
}

// SaveConfig 保存配置到文件
// 这个函数将当前的配置对象保存到config.json文件中
// 通常用于：
// 1. 初始化配置文件
// 2. 通过管理界面修改配置后保存
// 3. 备份当前配置
func SaveConfig(config *Config) error {
	// 将配置对象转换为格式化的JSON
	// json.MarshalIndent使JSON文件更易读（有缩进和换行）
	// 参数说明：
	//   config - 要序列化的对象
	//   ""     - 前缀，这里为空
	//   "  "   - 缩进，这里是两个空格
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		// JSON序列化失败，可能配置结构有问题
		return err
	}

	// 将JSON数据写入config.json文件
	// 参数说明：
	//   "config.json" - 文件名
	//   data          - 要写入的数据
	//   0644          - 文件权限：所有者可读写，其他人只读
	//                   6 = 110（二进制）= rw-（所有者权限）
	//                   4 = 100（二进制）= r--（组权限）
	//                   4 = 100（二进制）= r--（其他用户权限）
	return os.WriteFile("config.json", data, 0644)
}
