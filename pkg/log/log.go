package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"my_destributed_project/configs"
	"os"
	"path/filepath"
)

// Logger 全局日志变量
var Logger *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger() {
	// 读取日志配置
	logPath := configs.AppConfig.Server.LogPath
	logName := configs.AppConfig.Server.LogName

	// 处理 logPath 为空的情况
	if logPath == "" {
		fmt.Println("Log file path is empty! Please check configs.AppConfig.Server.LogPath.")
	}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	}

	// 确保日志目录存在
	logFullPath := filepath.Join(dir, logPath)
	if err := os.MkdirAll(logFullPath, 0755); err != nil {
		fmt.Println("Failed to create log directory:", err)
	}

	// 生成完整的日志文件路径
	logFile := filepath.Join(logFullPath, logName)

	// 配置日志切割
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile, // 日志文件路径
		MaxSize:    10,      // **最大日志文件大小 (MB)**，超过此大小会日志切割
		MaxAge:     30,      // **日志保留时间 (天)**，超过此时间的日志会删除
		MaxBackups: 5,       // **最多保留的历史日志文件数**
		Compress:   true,    // **是否启用 gzip 压缩**
	}

	// 日志格式配置
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建日志核心
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),    // 终端格式
		zapcore.AddSync(lumberjackLogger), // 日志切割写入
		zap.InfoLevel,                     // 记录 Info 及以上级别
	)

	// 创建日志器
	Logger = zap.New(core).Sugar()
	Logger.Info("Log initialized with rolling support")
}
