ifneq (,$(wildcard .env))
	include .env
	export
endif

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
	@go run cmd/server/main.go

