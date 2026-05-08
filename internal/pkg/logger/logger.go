// Package logger 基于 zap 的结构化日志封装。
// 约定:所有业务日志都带 trace_id 字段(由中间件注入到 ctx)。
package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置,对应 yaml log.*
type Config struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

type traceIDKey struct{}

// WithTraceID 把 trace_id 放到 context。
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, id)
}

// TraceIDFrom 从 context 提取 trace_id。
func TraceIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey{}).(string); ok {
		return v
	}
	return ""
}

var global *zap.Logger = zap.NewNop()

// Init 初始化全局 logger。
func Init(cfg Config) error {
	level, err := zapcore.ParseLevel(defaultIfEmpty(cfg.Level, "info"))
	if err != nil {
		return fmt.Errorf("parse log level: %w", err)
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.MessageKey = "msg"
	encCfg.LevelKey = "level"
	encCfg.CallerKey = "caller"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	var encoder zapcore.Encoder
	if defaultIfEmpty(cfg.Format, "json") == "console" {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	writer, err := openWriter(cfg)
	if err != nil {
		return fmt.Errorf("open log writer: %w", err)
	}

	core := zapcore.NewCore(encoder, writer, level)
	global = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

// L 返回不带 ctx 的全局 logger。
func L() *zap.Logger { return global }

// C 返回带 trace_id 的 logger(如果 ctx 里有的话)。
func C(ctx context.Context) *zap.Logger {
	if tid := TraceIDFrom(ctx); tid != "" {
		return global.With(zap.String("trace_id", tid))
	}
	return global
}

// Sync 刷盘。
func Sync() { _ = global.Sync() }

func openWriter(cfg Config) (zapcore.WriteSyncer, error) {
	switch defaultIfEmpty(cfg.Output, "stdout") {
	case "stdout":
		return zapcore.AddSync(os.Stdout), nil
	case "file":
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("log.file_path required when output=file")
		}
		f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		return zapcore.AddSync(f), nil
	default:
		return zapcore.AddSync(os.Stdout), nil
	}
}

func defaultIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
