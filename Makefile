APP_NAME := rss-aggregator
APP_PATH := ./cmd/api
ADDR ?= :5000

.PHONY: run build

run:
	@echo "Running $(APP_NAME)..."
	go run $(APP_PATH) -addr=$(ADDR)

build:
	@echo "Building $(APP_NAME)..."
	go build -o ./bin/$(APP_NAME) $(APP_PATH)