package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type Logger interface {
	Info(msg string, tags ...zap.Field)
	Error(msg string, err error, tags ...zap.Field)
}

type logger struct {
	log *zap.Logger
}

var (
	logInstance *logger
	once        sync.Once
)

func NewLogger() Logger {
	once.Do(func() {
		conf := zap.Config{
			OutputPaths: []string{"stdout"},
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:   "msg",
				LevelKey:     "level",
				TimeKey:      "time",
				EncodeLevel:  zapcore.LowercaseLevelEncoder,
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}

		log, err := conf.Build()
		if err != nil {
			panic(err)
		}

		logInstance = &logger{
			log: log,
		}

	})
	return logInstance
}

func (l *logger) Info(msg string, tags ...zap.Field) {
	l.log.Info(msg, tags...)
	_ = l.log.Sync()
}

func (l *logger) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Error(msg, tags...)
	_ = l.log.Sync()
}
