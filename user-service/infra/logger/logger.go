package logger

import (
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var Log *zap.Logger

func InitLogger() {
	once.Do(func() {
		var err error
		Log, err = zap.NewDevelopment()
		if err != nil {
			panic("No se pudo inicializar zap logger: " + err.Error())
		}
		Log.Info("✅ Logger Inicializado")
	})
}
