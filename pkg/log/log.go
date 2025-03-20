package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger 全局日志变量
var Logger *zap.SugaredLogger

// InitLogger 初始化日志
//func InitLogger() {
//	log, _ := zap.NewProduction()
//	Logger = log.Sugar()
//}

// InitLogger 初始化日志
func InitLogger(logPath string) {

	dir, getDirErr := os.Getwd()
	if nil != getDirErr {
		println("can not get current dir")
		return
	}
	println("current dir:", dir)

	// 创建一个日志文件的配置
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置日志文件的路径
	file, err := os.OpenFile(dir+logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	// 创建一个日志写入器
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	// 创建一个日志记录器
	log := zap.New(core)
	Logger = log.Sugar()
	println("日志文件创建成功:", dir+logPath)
}
