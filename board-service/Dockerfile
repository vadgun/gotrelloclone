# Usamos una imagen oficial de Golang
FROM golang:1.24 AS builder

# Establecemos el directorio de trabajo
WORKDIR /app

# Instalamos librdkafka
RUN apt-get update && apt-get install -y librdkafka-dev

# Copiamos el código del servicio
COPY . .

# Descargamos las dependencias
RUN go mod tidy

# Compilamos el binario
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o board-service .

# Segunda etapa: imagen debian-slim con soporte para librdkafka
FROM debian:12-slim

# Instalamos solo la biblioteca compartida de runtime
RUN apt-get update && apt-get install -y librdkafka1 && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copiamos el binario desde la fase de compilación
COPY --from=builder /app/board-service .

# Otorgamos permisos al servicio
RUN chmod +x board-service

# Exponemos el puerto
EXPOSE 8080

# Ejecutamos el servicio
CMD ["./board-service"]