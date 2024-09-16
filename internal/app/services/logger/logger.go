package logger

import (
	"go.uber.org/zap"
)

var LogSugar *zap.SugaredLogger

func NewLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return err
	}
	LogSugar = appLogger.Sugar()
	return nil
}
