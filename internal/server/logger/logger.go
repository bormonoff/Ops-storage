package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const spaceSeparator string = " "

// A sugar log needs for a func that isn't impact on performance
var MainLog *zap.SugaredLogger = zap.NewNop().Sugar()

// A strongly typed logger for a highload func
var HandlerLog *zap.Logger = zap.NewNop()

func genRootCfg() zap.Config {
	lvlEnc := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	enc := zapcore.EncoderConfig{
		LevelKey: "level",
		TimeKey:  "time",

		EncodeLevel:  lvlEnc,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,

		LineEnding:       zapcore.DefaultLineEnding,
		ConsoleSeparator: spaceSeparator,
	}

	rootCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     enc,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	return rootCfg
}

func initLogger() *zap.Logger {
	cfg := genRootCfg()

	cfg.EncoderConfig.MessageKey = "msg"
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Initialize() {
	MainLog = initLogger().Sugar()
	HandlerLog = initLogger()
}
