# Builder stage
FROM golang:1.22-alpine AS builder

# Set GOPROXY for stable downloads
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}

WORKDIR /app

# Copy source code first
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o babygram main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/babygram .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./babygram"]
