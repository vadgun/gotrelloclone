package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewDevelopment() // Puedes usar zap.NewDevelopment() en dev || zap.NewProduction en prod
	if err != nil {
		panic("No se pudo inicializar zap logger: " + err.Error())
	}

	Log.Info("✅ Logger Inicializado")
}
