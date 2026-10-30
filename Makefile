APP_NAME := rss-aggregator
APP_PATH := ./cmd/api
ADDR ?= :5000
DSN ?= postgres://user:password@localhost:5432/rssagg?sslmode=disable

.PHONY: run build

run:
	@echo "Running $(APP_NAME)..."
	go run $(APP_PATH) -addr=$(ADDR) -dsn=$(DSN)

build:
	@echo "Building $(APP_NAME)..."
	go build -o ./bin/$(APP_NAME) $(APP_PATH)