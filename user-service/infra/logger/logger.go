package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// Configura logrus para usar formato JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Configura el nivel de logging (puedes cambiarlo según el entorno)
	logrus.SetLevel(logrus.InfoLevel)

	// Configura la salida de los logs (por defecto es os.Stdout)
	logrus.SetOutput(os.Stdout)
}
