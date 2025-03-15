package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Port     string `mapstructure:"port"`
	SslMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

var AppConfig DatabaseConfig

// LoadConfig 读取 YAML 配置文件
func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.UnmarshalKey("database", &AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal database config: %w", err)
	}

	return nil
}
