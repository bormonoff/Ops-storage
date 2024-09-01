package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger = zap.NewNop().Sugar()

func Initialize() {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Log = logger.Sugar()
}
