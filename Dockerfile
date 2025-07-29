#FROM golang:1.21-alpine AS builder
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o stockconsolidation ./cmd/stockconsolidation

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/stockconsolidation .

# Create logs directory and set permissions
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/logs && \
    chown -R appuser:appgroup /app && \
    chmod -R 777 /app/logs

VOLUME ["/app/logs"]

USER appuser

EXPOSE 3000

CMD ["./stockconsolidation"]
