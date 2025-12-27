# Tools
GO_BIN := $(shell which go)
BUN_BIN := $(shell which bun)
FSWATCH_BIN := $(shell which fswatch)
PROJECT_NAME := $(shell basename $(CURDIR))
TEMPL_BIN := $(shell which templ)
ENTR_BIN := $(shell which entr)

# Project info
PROJECT_NAME := $(shell basename $(CURDIR))
STATIC_DIR := static

# Build configuration
BUILD_DIR := frontend/build
BACKEND_MAIN_GO := backend/main.go
FRONTEND_MAIN_GO := frontend/main.go
TEMPL_FILES := $(shell find . -type f -name "*.templ")
TEMPL_GO_FILES := $(TEMPL_FILES:.templ=_templ.go)

# Simple color scheme
CYAN := \033[36m
DIM := \033[2m
RESET := \033[0m

# Basic symbols
CHECK := ✓
ARROW := →

# URLs and resources
ALPINE_URL := https://cdn.jsdelivr.net/npm/alpinejs@latest/dist/cdn.min.js

# Simplified Tailwind config focused on templ files
TAILWIND_CONFIG := 'module.exports = {\n\
  content: ["./views/**/*.templ"],\n\
  theme: {\n\
    extend: {}\n\
  },\n\
  plugins: []\n\
}'

.PHONY: check-deps init create-dirs setup-go setup-tailwind download-alpine build serve watch clean help templ

check-deps:
	@printf "\n$(CYAN)Checking dependencies$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@test -n "$(GO_BIN)" || (printf "✗ Go not installed\n" && exit 1)
	@test -n "$(BUN_BIN)" || (printf "✗ Bun not installed\n" && exit 1)
	@test -n "$(FSWATCH_BIN)" || (printf "✗ fswatch not installed\n" && exit 1)
	@test -n "$(TEMPL_BIN)" || (printf "✗ templ not installed\n" && exit 1)
	@test -n "$(ENTR_BIN)" || (printf "✗ entr not installed\n" && exit 1)
	@printf "$(CHECK) All dependencies ready\n"

setup-tailwind:
	@printf "\n$(CYAN)Setting up Tailwind CSS$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@cd frontend && $(BUN_BIN) add tailwindcss@latest >/dev/null 2>&1
	@cd frontend && $(BUN_BIN) add -D @tailwindcss/cli>/dev/null 2>&1
	@cd frontend && echo $(TAILWIND_CONFIG) > tailwind.config.js
	@echo '@tailwind base;\n@tailwind components;\n@tailwind utilities;' > frontend/$(STATIC_DIR)/css/input.css
	@printf "$(CHECK) Tailwind CSS ready\n"

download-alpine:
	@printf "\n$(CYAN)Downloading Alpine.js$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@cd frontend && mkdir -p $(STATIC_DIR)/js
	@cd frontend && curl -s $(ALPINE_URL) -o $(STATIC_DIR)/js/alpine.min.js
	@printf "$(CHECK) Alpine.js ready\n"

build: check-deps
	@printf "\n$(CYAN)Building project$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@$(GO_BIN) build -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_GO)
	@printf "$(CHECK) Build complete\n"

templ:
	@printf "\n$(CYAN)Generating templ files$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@$(TEMPL_BIN) generate
	@printf "$(CHECK) Templ files generated\n"

css:
	@printf "\n$(CYAN)Generating CSS$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@cd frontend && $(BUN_BIN) tailwindcss -i $(STATIC_DIR)/css/input.css -o $(STATIC_DIR)/css/styles.css --minify
	@printf "$(CHECK) CSS generated\n"

backend_serve: 
	@printf "\n$(CYAN)Starting backend server$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@$(GO_BIN) run $(BACKEND_MAIN_GO)

frontend_serve: templ css
	@printf "\n$(CYAN)Starting frontend server$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@cd frontend && $(GO_BIN) run main.go serve 

watch:
	@printf "\n$(CYAN)Watching for changes$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@printf "\n$(CYAN)Rebuilding...$(RESET)\n"
	@find frontend/views -type f -name "*.templ" | entr -r make frontend_serve

clean:
	@printf "\n$(CYAN)Cleaning build files$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@rm -rf $(BUILD_DIR)
	@find . -type f -name "*_templ.go" -delete
	@printf "$(CHECK) Clean complete\n"

init: check-deps setup-tailwind download-alpine
	@printf "\n$(CHECK) Project setup complete!\n"

test:
	@printf "\n$(CYAN)Running tests$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@$(GO_BIN) test -v -count=1 -cover ./...

source:
	@printf "\n$(CYAN)Loading environment variables (dotenv)$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@export $(cat .env | xargs)
	@printf "\n$(CHECK) Environment variables loaded\n"

sync:
	@printf "\n$(CYAN)Running browser-sync$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@browser-sync start --proxy "localhost:8090" --files "views/**/*.templ"

help:
	@printf "\n$(CYAN)Available commands$(RESET)\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@printf "  make init$(RESET)          $(ARROW) Initialize project\n"
	@printf "  make build$(RESET)         $(ARROW) Build project\n"
	@printf "  make css$(RESET)           $(ARROW) Generate CSS\n"
	@printf "  make serve$(RESET)         $(ARROW) Start server\n"
	@printf "  make watch$(RESET)         $(ARROW) Watch for changes\n"
	@printf "  make clean$(RESET)         $(ARROW) Clean build files\n"
	@printf "  make source$(RESET)        $(ARROW) Load environment variables\n"
	@printf "  make sync$(RESET)          $(ARROW) run browser-sync\n"
	@printf "  make test$(RESET)          $(ARROW) Run tests\n\n"
