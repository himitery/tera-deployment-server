package logger

import (
	"fmt"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"tera/deployment/pkg/config"
)

var (
	log *zap.Logger
)

func Init(conf *config.Config) *zap.Logger {
	level := lo.Switch[string, zap.AtomicLevel](conf.Logging.Level).
		Case("debug", zap.NewAtomicLevelAt(zap.DebugLevel)).
		Case("info", zap.NewAtomicLevelAt(zap.InfoLevel)).
		Case("warn", zap.NewAtomicLevelAt(zap.WarnLevel)).
		Case("error", zap.NewAtomicLevelAt(zap.ErrorLevel)).
		Default(zap.NewAtomicLevelAt(zap.InfoLevel))

	encoding := lo.Switch[string, string](conf.Profile).
		Case("development", "console").
		Default("json")

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	base, err := zap.Config{
		Level:            level,
		Development:      conf.Profile == "development",
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if err != nil {
		fmt.Println(err)
		panic("can't initialize logger")
	}

	log = base.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	return log
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}
