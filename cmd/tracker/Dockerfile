# Etapa 1: build
FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar el binario (forzar Linux + estático)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tracker ./tracker

# Etapa 2: imagen final mínima
FROM alpine:3.20

RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /app/tracker .

USER appuser
CMD ["./tracker"]
