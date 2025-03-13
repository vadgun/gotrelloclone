# Usamos una imagen oficial de Golang
FROM golang:1.24 AS builder

# Establecemos el directorio de trabajo
WORKDIR /app

# Copiamos el código del servicio
COPY . .

# Copiamos el archivo .env del servicio
COPY .env .

# Descargamos las dependencias
RUN go mod tidy

# Compilamos el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o board-service .

# Segunda etapa: imagen ligera de ejecución
FROM alpine:latest  

WORKDIR /root/

# Copiamos el binario desde la fase de compilación
COPY --from=builder /app/board-service .

# Otorgamos permisos al servicio
RUN chmod +x board-service

# Exponemos el puerto
EXPOSE 8080

# Ejecutamos el servicio
CMD ["./board-service"]
