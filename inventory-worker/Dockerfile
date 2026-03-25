FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY inventory-worker/go.mod inventory-worker/go.sum ./inventory-worker/

WORKDIR /app/inventory-worker

RUN go mod download

COPY inventory-worker ./ 

RUN go build -o inventory-worker ./cmd

FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/inventory-worker/inventory-worker .

RUN chown appuser:appuser /app/inventory-worker

USER appuser

CMD ["./inventory-worker"]