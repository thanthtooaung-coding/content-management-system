package logger

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
}

func NewLogger(params Params) *logrus.Logger {
	logger := logrus.New()

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	logger.SetOutput(lumberjackLogger)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing log file")
			return lumberjackLogger.Close()
		},
	})

	return logger
}
