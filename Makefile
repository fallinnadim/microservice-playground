.PHONY: run-dev

run-dev:
	@go run ./order-service/cmd/api & \
	go run ./payment-service/cmd/api & \
	go run ./inventory-worker/cmd & \
	wait

run-up:
	@docker-compose up -d --build

run-down:
	@docker-compose down
