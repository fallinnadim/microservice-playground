.PHONY: run-dev

run-dev:
	@go run ./order-service/cmd/api

order-service-build:
	@docker build -f order-service/Dockerfile -t order-service .

order-service-run:
	@docker run -d -p 8080:8080 --env-file order-service/.env --name order-service order-service:latest
