# Variables
PYTHON = python3
PIP = pip
PYINSTALLER = pyinstaller
PROJECT_NAME = tinymonitor
ENTRY_POINT = src/tinymonitor/__main__.py
DIST_DIR = dist
BUILD_DIR = build

.PHONY: help install test build clean run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies in editable mode
	$(PIP) install -e .
	$(PIP) install pyinstaller

test: ## Run unit tests
	$(PYTHON) -m unittest discover -s tests -p "test_*.py"

build: clean ## Build the standalone binary using PyInstaller
	$(PYINSTALLER) --name $(PROJECT_NAME) --onefile --paths src $(ENTRY_POINT)
	@echo "✅ Build complete! Binary: $(DIST_DIR)/$(PROJECT_NAME)"

clean: ## Clean up build artifacts and cache
	rm -rf $(BUILD_DIR) $(DIST_DIR) $(PROJECT_NAME).spec
	rm -rf src/*.egg-info
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete

run: ## Run the application from source
	$(PYTHON) -m tinymonitor -c config.json
