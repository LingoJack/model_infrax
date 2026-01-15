package logger

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *ILogger
)

func init() {
	fmt.Println("logger init...")
	logger, _ = NewLogger(&ILoggerConfig{
		LogPath:    "./logs/robot_client.log",
		MaxAge:     30,
		MaxSize:    100,
		MaxBackups: 10,
		Compress:   false,
		Level:      "info",
	})
}

type ILogger struct {
	Core   *zap.Logger
	Sugar  *zap.SugaredLogger
	writer *RollWriter
	*zap.SugaredLogger
}

func Logger() *ILogger {
	return logger
}

// Info 记录 Info 级别日志，会显示真实调用位置
func Info(args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Info(args...)
}

// Debug 记录 Debug 级别日志，会显示真实调用位置
func Debug(args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Debug(args...)
}

// Error 记录 Error 级别日志，会显示真实调用位置
func Error(args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Error(args...)
}

// Warn 记录 Warn 级别日志，会显示真实调用位置
func Warn(args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Warn(args...)
}

// Fatal 记录 Fatal 级别日志，会显示真实调用位置
func Fatal(args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatal(args...)
}

// Infof 记录格式化的 Info 级别日志，会显示真实调用位置
func Infof(format string, args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Infof(format, args...)
}

// Debugf 记录格式化的 Debug 级别日志，会显示真实调用位置
func Debugf(format string, args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugf(format, args...)
}

// Errorf 记录格式化的 Error 级别日志，会显示真实调用位置
func Errorf(format string, args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorf(format, args...)
}

// Warnf 记录格式化的 Warn 级别日志，会显示真实调用位置
func Warnf(format string, args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Warnf(format, args...)
}

// Fatalf 记录格式化的 Fatal 级别日志，会显示真实调用位置
func Fatalf(format string, args ...interface{}) {
	logger.Core.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalf(format, args...)
}

func NewLogger(cfg *ILoggerConfig) (*ILogger, error) {
	writer, err := NewRollWriter(cfg.LogPath, []Option{
		WithMaxAge(cfg.MaxAge),         // 日志文件最大保留天数
		WithMaxSize(cfg.MaxSize),       // 单个日志文件最大大小(MB)
		WithCompress(cfg.Compress),     // 是否压缩归档日志
		WithMaxBackups(cfg.MaxBackups), // 最大保留的归档日志文件数
	})
	if err != nil {
		log.Printf("create rollwriter failed: %v", err)
		return nil, err
	}

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	// 配置 zap encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 大写日志级别
		EncodeTime:     customTimeEncoder,             // 自定义人类易读时间格式
		EncodeDuration: zapcore.StringDurationEncoder, // 字符串持续时间编码
		EncodeCaller:   zapcore.FullCallerEncoder,     // 完整路径调用者信息（包含方法名）
	}

	// 解析日志级别
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	)

	// 创建 log，添加调用者信息，错误级别添加堆栈跟踪
	coreLogger := zap.New(core, zap.AddCaller())
	sugarLogger := coreLogger.Sugar()

	return &ILogger{
		Core:          coreLogger,
		Sugar:         sugarLogger,
		writer:        writer,
		SugaredLogger: sugarLogger,
	}, nil
}
