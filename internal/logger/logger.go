package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

func InitDefault(logLevel string) error {
	atomicLevel, err := getAtomicLevel(logLevel)
	if err != nil {
		return err
	}

	Init(getCore(atomicLevel))

	return nil
}

func Logger() *zap.Logger {
	return globalLogger
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}
func Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	globalLogger.Info(msg)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	globalLogger.Warn(msg)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	globalLogger.Error(msg)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	globalLogger.Fatal(msg)
}

func WithOptions(options ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(options...)
}
