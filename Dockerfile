FROM golang:1.22-alpine AS builder
WORKDIR /app

# Копируем только go.mod
COPY go.mod ./

# Обновляем зависимости и скачиваем их
RUN go mod tidy
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o babygram .

# Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/babygram .
COPY migrations ./migrations
EXPOSE 8080
CMD ["./babygram"]
