#!/bin/bash

echo "📦 Reorganizando estructura del user-service..."

SERVICE_DIR="./user-service"
INFRA_DIR="$SERVICE_DIR/infra"

# Crear estructura de carpetas
mkdir -p $INFRA_DIR/config
mkdir -p $INFRA_DIR/kafka
mkdir -p $INFRA_DIR/logger
mkdir -p $INFRA_DIR/metrics

# Mover archivos usando git mv si está disponible
if command -v git &> /dev/null; then
  echo "✅ Usando git mv para preservar historial..."

  git mv $SERVICE_DIR/config/user_config.go $INFRA_DIR/config/config.go
  git mv $SERVICE_DIR/kafka/producer.go $INFRA_DIR/kafka/producer.go
  git mv $SERVICE_DIR/logger/logger.go $INFRA_DIR/logger/logger.go
  git mv $SERVICE_DIR/metrics/metrics.go $INFRA_DIR/metrics/metrics.go

else
  echo "⚠️ git no encontrado, usando mv normal (historial se perderá)..."

  mv $SERVICE_DIR/config/user_config.go $INFRA_DIR/config/config.go
  mv $SERVICE_DIR/kafka/producer.go $INFRA_DIR/kafka/producer.go
  mv $SERVICE_DIR/logger/logger.go $INFRA_DIR/logger/logger.go
  mv $SERVICE_DIR/metrics/metrics.go $INFRA_DIR/metrics/metrics.go
fi

# Eliminar carpetas vacías (opcional)
rmdir $SERVICE_DIR/config 2> /dev/null
rmdir $SERVICE_DIR/kafka 2> /dev/null
rmdir $SERVICE_DIR/logger 2> /dev/null
rmdir $SERVICE_DIR/metrics 2> /dev/null

echo "✅ user-service reorganizado con la carpeta infra 🚀"
