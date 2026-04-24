package logger

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Init(env string) {
	once.Do(func() {
		var err error
		if env == "production" {
			instance, err = zap.NewProduction()
		} else {
			instance, err = zap.NewDevelopment()
		}
		if err != nil {
			panic(err)
		}
	})
}

func Get() *zap.Logger {
	if instance == nil {
		panic("logger não inicializado, chame Init() primeiro")
	}
	return instance
}

func Sync() {
	if instance != nil {
		instance.Sync()
	}
}

func LogError(message string) error {
	err := errors.New(message)
	Get().Error("datagram error", zap.Error(err))
	return err
}
