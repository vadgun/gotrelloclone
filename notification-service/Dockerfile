# Usamos una imagen oficial de Golang
FROM golang:1.23 AS builder

# Establecemos el directorio de trabajo
WORKDIR /app

# Copiamos el código del servicio
COPY . .

# Descargamos las dependencias
RUN go mod tidy

# Compilamos el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notification-service .

# Segunda etapa: imagen ligera de ejecución
FROM alpine:latest  

WORKDIR /root/

# Copiamos el binario desde la fase de compilación
COPY --from=builder /app/notification-service .

# Otorgamos permisos al servicio
RUN chmod +x notification-service

# Exponemos el puerto
EXPOSE 8080

# Ejecutamos el servicio
CMD ["./notification-service"]
