export PROJECT_ROOT=$(shell pwd)

chat-run:
	@go run ${PROJECT_ROOT}/cmd/chat/main.go