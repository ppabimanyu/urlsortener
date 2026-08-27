.PHONY: help setup dev db-up db-down db-logs db-reset be-run be-build be-test be-tidy fe-install fe-dev fe-build fe-typecheck docker-build docker-up docker-down docker-logs clean

# Default target
.DEFAULT_GOAL := help

# Colors for terminal output
COLOR_RESET   := \033[0m
COLOR_INFO    := \033[36m
COLOR_SUCCESS := \033[32m
COLOR_WARNING := \033[33m
COLOR_HEADER  := \033[1;34m

##@ 📖 General
help: ## Show this help message
	@echo ""
	@printf "${COLOR_HEADER}SortLink - Project Commands${COLOR_RESET}\n"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make ${COLOR_INFO}<target>${COLOR_RESET}\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  ${COLOR_INFO}%-18s${COLOR_RESET} %s\n", $$1, $$2 } /^##@/ { printf "\n${COLOR_HEADER}%s${COLOR_RESET}\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

setup: ## Initial project setup (create .env, install dependencies, start db)
	@printf "${COLOR_INFO}==> Setting up environment files...${COLOR_RESET}\n"
	@test -f backend/.env || cp backend/.env.example backend/.env
	@test -f frontend/.env || cp frontend/.env.example frontend/.env
	@printf "${COLOR_INFO}==> Installing frontend dependencies...${COLOR_RESET}\n"
	@cd frontend && npm install
	@printf "${COLOR_INFO}==> Starting database container...${COLOR_RESET}\n"
	@docker compose up -d postgres
	@printf "${COLOR_SUCCESS}✓ Setup complete! Run 'make dev' to start all services.${COLOR_RESET}\n"

dev: db-up ## Run both Backend and Frontend concurrently
	@printf "${COLOR_INFO}==> Starting Backend and Frontend concurrently... (Press Ctrl+C to stop)${COLOR_RESET}\n"
	@trap 'kill 0' SIGINT SIGTERM EXIT; \
	(cd backend && go run ./cmd/server) & \
	(cd frontend && npm run dev) & \
	wait

##@ 🗄️ Database
db-up: ## Start PostgreSQL database container
	@printf "${COLOR_INFO}==> Starting PostgreSQL container...${COLOR_RESET}\n"
	@docker compose up -d postgres

db-down: ## Stop PostgreSQL database container
	@printf "${COLOR_INFO}==> Stopping PostgreSQL container...${COLOR_RESET}\n"
	@docker compose stop postgres

db-logs: ## View PostgreSQL logs
	@docker compose logs -f postgres

db-reset: ## Reset PostgreSQL database and volumes
	@printf "${COLOR_WARNING}==> Resetting database and volumes...${COLOR_RESET}\n"
	@docker compose down -v
	@docker compose up -d postgres

##@ ⚙️ Backend (Go)
be-run: db-up ## Run Go backend server (http://localhost:8080)
	@printf "${COLOR_INFO}==> Running Backend server...${COLOR_RESET}\n"
	@cd backend && go run ./cmd/server

be-build: ## Build backend binary
	@printf "${COLOR_INFO}==> Building Backend binary...${COLOR_RESET}\n"
	@mkdir -p backend/bin
	@cd backend && go build -o bin/server ./cmd/server
	@printf "${COLOR_SUCCESS}✓ Binary built at backend/bin/server${COLOR_RESET}\n"

be-test: ## Run backend unit tests
	@printf "${COLOR_INFO}==> Running Backend tests...${COLOR_RESET}\n"
	@cd backend && go test -v ./...

be-tidy: ## Tidy and download Go modules
	@printf "${COLOR_INFO}==> Tidying Go modules...${COLOR_RESET}\n"
	@cd backend && go mod tidy

##@ 🎨 Frontend (Vue 3)
fe-install: ## Install frontend npm dependencies
	@printf "${COLOR_INFO}==> Installing Frontend dependencies...${COLOR_RESET}\n"
	@cd frontend && npm install

fe-dev: ## Run frontend development server (http://localhost:5173)
	@printf "${COLOR_INFO}==> Starting Frontend Vite dev server...${COLOR_RESET}\n"
	@cd frontend && npm run dev

fe-build: ## Build frontend for production
	@printf "${COLOR_INFO}==> Building Frontend for production...${COLOR_RESET}\n"
	@cd frontend && npm run build

fe-typecheck: ## Run frontend TypeScript type checking
	@printf "${COLOR_INFO}==> Checking Frontend TypeScript types...${COLOR_RESET}\n"
	@cd frontend && npx vue-tsc -b

##@ 🐳 Docker
docker-build: ## Build unified Docker image (Backend + Frontend)
	@printf "${COLOR_INFO}==> Building Docker container...${COLOR_RESET}\n"
	@docker compose build

docker-up: ## Start all services in Docker (PostgreSQL & App on port 8080)
	@printf "${COLOR_INFO}==> Starting services with Docker Compose...${COLOR_RESET}\n"
	@docker compose up -d

docker-down: ## Stop all Docker containers
	@printf "${COLOR_INFO}==> Stopping all Docker containers...${COLOR_RESET}\n"
	@docker compose down

docker-logs: ## View logs of all Docker services
	@docker compose logs -f

##@ 🧹 Maintenance
clean: ## Clean build artifacts and binaries
	@printf "${COLOR_INFO}==> Cleaning build artifacts...${COLOR_RESET}\n"
	@rm -rf backend/bin
	@rm -rf frontend/dist
	@printf "${COLOR_SUCCESS}✓ Cleaned!${COLOR_RESET}\n"
