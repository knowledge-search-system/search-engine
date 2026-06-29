package logger

import (
	"github.com/knowledge-search-system/search-engine/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Module("logger",
	fx.Provide(New),
)

func New(cfg *config.Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.Set(cfg.Logger.Level); err != nil {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapCfg.Build()
}
