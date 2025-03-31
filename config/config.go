package config

import (
	"github.com/spf13/viper"
	"log"
)

// ServerConfig 定义服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// AIProxyConfig 定义 AI Proxy 配置
type AIProxyConfig struct {
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

// DatabaseConfig 定义数据库配置
type DatabaseConfig struct {
	Type           string `mapstructure:"type"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	DatabaseName   string `mapstructure:"database_name"`
	ConnectionPool struct {
		MaxConnections int `mapstructure:"max_connections"`
		MinConnections int `mapstructure:"min_connections"`
		IdleTimeout    int `mapstructure:"idle_timeout"`
	} `mapstructure:"connection_pool"`
}

// LoggingConfig 定义日志配置
type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// JWTConfig 定义 JWT 配置
type JWTConfig struct {
	Secret string `mapstructure:"jwt_sec"` // 对应配置文件中的 jwt_sec
}

// Config 包含所有的配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	AIProxy  AIProxyConfig  `mapstructure:"ai_proxy"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	JWT      JWTConfig      `mapstructure:"jwt"` // 添加 JWT 配置
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 设置配置文件路径
	viper.AddConfigPath("./config") // 配置文件所在目录
	viper.SetConfigName("config")   // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")     // 配置文件类型

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// InitViper 提供直接使用 viper 的方法
func InitViper() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}
}

func GetAPIURL() string {
	return viper.GetString("ai_proxy.url")
}

func GetAPIKey() string {
	return viper.GetString("ai_proxy.api_key")
}
