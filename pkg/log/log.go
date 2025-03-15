package log

import (
	"go.uber.org/zap"
)

// Logger 全局日志变量
var Logger *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger() {
	log, _ := zap.NewProduction()
	Logger = log.Sugar()
}
