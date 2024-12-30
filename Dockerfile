# Etapa de construcción
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copia el binario compilado y el archivo .env
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
