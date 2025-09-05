.PHONY: dev build run clean kill

# Development mode with hot reload and verbose logging
dev: kill
	@echo "Starting in development mode..."
	@mkdir -p build/pb_public/email_templates
	@mkdir -p build/pb_public/bot_templates
	@cp -r src/static/* build/pb_public/ 2>/dev/null || true
	@cp .env build/.env 2>/dev/null || true
	@cd build && . ./.env && go run ../src/main.go serve --dev --http=0.0.0.0:$${PORT:-8080}

# Kill any running instance
kill:
	@echo "Killing existing processes..."
	@-pkill -f "disciplo" 2>/dev/null || true
	@-pkill -f "go run.*main.go" 2>/dev/null || true
	@-pkill -f "main.go" 2>/dev/null || true
	@-lsof -ti:8080 | xargs -r kill -9 2>/dev/null || true
	@-lsof -ti:8081 | xargs -r kill -9 2>/dev/null || true
	@-lsof -ti:$${PORT:-8080} | xargs -r kill -9 2>/dev/null || true
	@sleep 2

# Build for production
build:
	@echo "Building for production..."
	@mkdir -p build/pb_public/email_templates
	@mkdir -p build/pb_public/bot_templates
	@go build -o build/disciplo src/main.go
	@cp -r src/static/* build/pb_public/ 2>/dev/null || true
	@cp .env build/.env 2>/dev/null || true
	@echo "Build complete. Binary at build/disciplo"

# Run the built application
run:
	@cd build && . ./.env && ./disciplo serve --dev --http=0.0.0.0:$${PORT:-8080}

# Clean build artifacts and database
clean:
	@echo "Cleaning build artifacts and database..."
	@rm -rf build/
	@echo "Clean complete"
