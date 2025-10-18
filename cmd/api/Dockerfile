FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar el binario con nombre fijo en /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/api-bin ./api

# Imagen final
FROM alpine:3.20
RUN adduser -D appuser
WORKDIR /app

# Copiamos explícitamente el binario desde la etapa builder
COPY --from=builder /app/api-bin /app/api-bin

USER appuser
CMD ["/app/api-bin"]
