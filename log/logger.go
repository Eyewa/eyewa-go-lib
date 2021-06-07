package log

import (
	"fmt"
	"os"

	"github.com/ory/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger    *zap.Logger
	logLevels map[string]zapcore.Level
)

func init() {
	logLevels = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"fatal": zap.FatalLevel,
		"error": zap.ErrorLevel,
	}
}

// SetLogLevel sets the log level detected from the env - LOG_LEVEL
// If the LOG_LEVEL is not set/found, the log level will default to INFO
func SetLogLevel() {
	viper.AutomaticEnv()
	level := os.Getenv("LOG_LEVEL")
	fmt.Println("I am LEVEL", level)

	lvl, ok := logLevels[level]
	if !ok {
		bootstrap(zapcore.InfoLevel)
		return
	}

	bootstrap(lvl)
}

func bootstrap(level zapcore.Level) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	logger = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevelAt(level),
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

// DebugWithTracing log debug entry with a trace id
func DebugWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Debug(message, fields...)
}

// Debug log entry at debug level
func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

// InfoWithTraceID log info entry with a trace id
func InfoWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Info(message, fields...)
}

// Info log entry at info level
func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

// WarnWithTraceID log warning entry with a trace id
func WarnWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Warn(message, fields...)
}

// Warn log entry at warning level
func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

// ErrorWithTraceID log error entry with a trace id
func ErrorWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Error(message, fields...)
}

// Error log entryat error level
func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

// FatalWithTraceID log fatal entry with a trace id
func FatalWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Fatal(message, fields...)
}

// Fatal log entry at fatal level
func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

// PanicWithTraceID log panic entry with a trace id
func PanicWithTraceID(traceID, message string, fields ...zap.Field) {
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	Panic(message, fields...)
}

// Panic log entry at panic level
func Panic(message string, fields ...zap.Field) {
	logger.Panic(message, fields...)
}
