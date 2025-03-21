package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

// Database 数据库配置结构体
type Database struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Port     string `mapstructure:"port"`
	SslMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Port    string `mapstructure:"port"`
	LogPath string `mapstructure:"log_path"`
	LogName string `mapstructure:"log_name"`
}

// Config 总体配置
type Config struct {
	Database Database     `mapstructure:"database"`
	Server   ServerConfig `mapstructure:"server"`
}

// AppConfig 全局配置变量
var AppConfig Config

// LoadConfig 读取 YAML 配置文件
func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
