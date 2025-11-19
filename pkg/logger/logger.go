package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	Log  *zap.Logger
	once sync.Once
)

func NewZapLogger() {
	once.Do(func() {
		var err error
		Log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	})
}

func init() {
	NewZapLogger()
}
