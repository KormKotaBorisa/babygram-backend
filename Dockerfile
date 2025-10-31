FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o babygram .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/babygram .
COPY migrations ./migrations
EXPOSE 8080
CMD ["./babygram"]
