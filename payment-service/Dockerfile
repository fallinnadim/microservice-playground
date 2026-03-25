FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY payment-service/go.mod payment-service/go.sum ./payment-service/

WORKDIR /app/payment-service

RUN go mod download

COPY payment-service ./ 

RUN go build -o payment-service ./cmd/api

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/payment-service/payment-service .

RUN chown appuser:appuser /app/payment-service

USER appuser

EXPOSE 8080

CMD ["./payment-service"]