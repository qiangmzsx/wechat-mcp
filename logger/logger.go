// Package logger 提供全局日志实例和便捷方法
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// 全局 logger 实例
	globalLogger *zap.Logger
)

// Init 初始化全局 logger
func Init(level, format string) error {
	zapLevel := zapcore.InfoLevel
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.AddSync(nil)),
		zapLevel,
	)

	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

// InitWithConfig 使用配置初始化全局 logger
func InitWithConfig(cfg *zap.Config) error {
	var err error
	globalLogger, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

// InitWithLogger 使用已有的 logger 初始化全局实例
func InitWithLogger(logger *zap.Logger) {
	globalLogger = logger
}

// Get 获取全局 logger 实例
func Get() *zap.Logger {
	if globalLogger == nil {
		// 返回 Nop logger 防止 nil panic
		return zap.NewNop()
	}
	return globalLogger
}

// Sync 同步日志缓冲区
func Sync() {
	if globalLogger != nil {
		globalLogger.Sync()
	}
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	Get().Sugar().Debugf(template, args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	Get().Sugar().Infof(template, args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	Get().Sugar().Warnf(template, args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	Get().Sugar().Errorf(template, args...)
}
