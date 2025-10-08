# ===================================================
# Go REST API with Swagger, Air & Docker
# ===================================================

APP_NAME := rest-api
PORT := 8089

# === Setup tools ===
.PHONY: tools
tools:
	@echo "→ Installing dev tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/air-verse/air@latest

# === Generate Swagger ===
.PHONY: swag
swag:
	@echo "→ Generating Swagger docs..."
	@swag init -g main.go -q

# === Run with hot reload ===
.PHONY: dev
dev: swag
	@echo "→ Starting dev server on port $(PORT) with Air..."
	@air

# === Run without hot reload ===
.PHONY: run
run: swag
	@echo "→ Running Go app..."
	@go run .

# === Build binary ===
.PHONY: build
build:
	@echo "→ Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) .

# === Docker build ===
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(APP_NAME):latest .

# === Docker run ===
.PHONY: docker-run
docker-run:
	@echo "🐳 Running Docker container on port $(PORT)..."
	@docker run --rm \
		-p $(PORT):8089 \
		--env-file .env \
		-v $(PWD)/.env:/app/.env:ro \
		$(APP_NAME):latest

# === Clean project ===
.PHONY: clean
clean:
	@echo "→ Cleaning build artifacts..."
	@rm -rf bin tmp

# === Show available commands ===
.PHONY: help
help:
	@echo ""
	@echo "🛠  Available commands:"
	@echo "  make tools         - install swag & air"
	@echo "  make swag          - generate Swagger docs"
	@echo "  make dev           - run with Air (hot reload)"
	@echo "  make run           - run without Air"
	@echo "  make build         - compile Go binary"
	@echo "  make docker-build  - build Docker image"
	@echo "  make docker-run    - run container on port $(PORT)"
	@echo "  make clean         - remove build files"
	@echo ""
