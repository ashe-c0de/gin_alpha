package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my_destributed_project/configs"
	"my_destributed_project/pkg/log"
)

// ConnectDatabase 连接数据库
func ConnectDatabase() *gorm.DB {
	// 读取数据库连接配置
	cfg := configs.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Logger.Error("Failed to connect to database:", err)
	}
	return db
}
