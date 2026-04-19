include .env
export

export PROJECT_ROOT=$(shell pwd)

up-db:
	@docker compose up -d postgres-db

forward-port:
	@docker compose up -d port-forwarder

run-chat:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export JWT_PRIVATE_PATH=${PROJECT_ROOT}/certs/app.rsa && \
	export JWT_PUBLIC_PATH=${PROJECT_ROOT}/certs/app.rsa.pub && \
	export POSTGRES_HOST=172.17.0.1 && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/chat/main.go

deploy-chat:
	@docker compose up -d --build chat-app

create-migration:
	@if [ -z "$(seq)" ]; then \
		echo "seq not set!"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-go-migration \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

up-migration:
	@make migrate-action action=up

down-migration:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action not set!"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-go-migration \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

clean-logs:
	@read -p "Clean all log files? Logs will be lost. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		sudo rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Logs was cleaned"; \
	else \
		echo "Clean canceled"; \
	fi