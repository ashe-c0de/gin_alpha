package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"my_destributed_project/configs"
	"os"
	"path/filepath"
	"time"
)

// Logger 全局日志变量
var Logger *zap.Logger

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
		MaxSize:    10,      // 最大日志文件大小 (MB)，超过此大小会日志切割
		MaxAge:     30,      // 日志保留时间 (天)，超过此时间的日志会删除
		MaxBackups: 5,       // 最多保留的历史日志文件数
		Compress:   true,    // 是否启用 gzip 压缩
	}

	// 创建标准输出的核心
	stdoutCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(getEncoderConfig()), // 终端友好格式
		zapcore.AddSync(os.Stdout),                    // 输出到标准输出
		zap.InfoLevel,                                 // 记录 Info 及以上级别
	)

	// 创建文件输出的核心
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(getEncoderConfig()), // JSON格式
		zapcore.AddSync(lumberjackLogger),          // 日志切割写入
		zap.InfoLevel,                              // 记录 Info 及以上级别
	)

	// 使用Tee Core组合多个输出
	core := zapcore.NewTee(stdoutCore, fileCore)

	// 创建日志器，添加调用者信息
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger.Info("Log initialized with rolling support")
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig() zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return cfg
}

// timeEncoder 自定义时间编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000-0700"))
}

// Debug 为了方便使用，添加一些辅助方法
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
