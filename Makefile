.PHONY: build
build:
	@echo "build gophermart"
	@go build -o cmd/gophermart/gophermart cmd/gophermart/main.go

.PHONY: run
run:
	@echo "run gophermart"
	@go run cmd/gophermart/main.go \
		-a :8080 \
		-d "postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable" \
		-r http://localhost:8081

.PHONY: test
test:
	@echo "test"
	@go test ./... -v

.PHONY: cover
cover:
	@echo "coverage"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep "^total:"

.PHONY: generate
generate:
	@echo "generate"
	@go generate ./...

.PHONY: infra-up
infra-up:
	@echo "infra-up"
	@docker compose -f deployments/docker-compose.yaml up -d

.PHONY: infra-down
infra-down:
	@echo "infra-down"
	@docker compose -f deployments/docker-compose.yaml down
