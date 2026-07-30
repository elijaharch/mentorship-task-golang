ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: all
all: lint test

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: db-up
db-up:
	docker compose up -d database

.PHONY: db-down
db-down:
	docker compose down database

.PHONY: migrate-up
migrate-up:
	goose up

.PHONY: migrate-down
migrate-down:
	goose down

.PHONY: migrate-status
migrate-status:
	goose status

.PHONY: migrate-new
migrate-new:
ifndef name
	$(error "Usage: make migrate-new name=create_users_table")
endif
	goose -dir $(MIGRATIONS_DIR) -s create $(name) sql

.PHONY: dev
dev:
	air

.PHONY: build
build:
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_BIN) ./cmd/$(APP_NAME)
	@echo "Built: $(OUT_BIN)"

.PHONY: clean
clean:
	@rm -rf $(OUT_DIR) tmp

.PHONY: test
test:
	@go test ./...

.PHONY: test-integration
test-integration:
	@go test -tags integration ./test/integration/...

.PHONY: run
run:
	@go run cmd/server/main.go

.PHONY: fmt
fmt:
	go fmt ./...
	yamlfmt .

.PHONY: lint
lint: tidy tools fmt security
	golangci-lint run ./...
	go vet ./...

.PHONY: security
security:
	govulncheck ./...

.PHONY: tools
tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest
	go install github.com/grpc-ecosystem/grpc-health-probe@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
