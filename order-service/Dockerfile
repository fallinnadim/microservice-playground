FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY order-service/go.mod order-service/go.sum ./order-service/

WORKDIR /app/order-service

RUN go mod download

COPY order-service ./ 

RUN go build -o order-service ./cmd/api

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/order-service/order-service .

RUN chown appuser:appuser /app/order-service

USER appuser

EXPOSE 8080

CMD ["./order-service"]