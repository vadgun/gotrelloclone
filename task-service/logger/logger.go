package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction() // Puedes usar zap.NewDevelopment() en dev
	if err != nil {
		panic("No se pudo inicializar zap logger: " + err.Error())
	}
}
