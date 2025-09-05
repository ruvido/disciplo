.PHONY: dev dev-email build run clean kill

# Development mode with hot reload and verbose logging
dev:
	@echo "Starting in development mode (emails are logged, not sent)..."
	@mkdir -p build/pb_public/email_templates
	@mkdir -p build/pb_public/bot_templates
	@cp -r src/static/* build/pb_public/ 2>/dev/null || true
	@cp .env build/.env 2>/dev/null || true
	@cd build && . ./.env && go run ../src/main.go serve --dev --http=0.0.0.0:$${PORT:-8080}

# Development mode with real email sending
dev-email:
	@echo "Starting in development mode with REAL email sending..."
	@mkdir -p build/pb_public/email_templates
	@mkdir -p build/pb_public/bot_templates
	@cp -r src/static/* build/pb_public/ 2>/dev/null || true
	@cp .env build/.env 2>/dev/null || true
	@cd build && . ./.env && go run ../src/main.go serve --http=0.0.0.0:$${PORT:-8080}

# Kill any running instance - BRUTAL MODE
kill:
	@echo "🔥 KILLING ALL PROCESSES - BRUTAL MODE..."
	@# Kill all disciplo processes
	@pkill -9 -f "disciplo" 2>/dev/null || true
	@pkill -9 -f "go run.*main.go" 2>/dev/null || true
	@pkill -9 -f "go-build.*main" 2>/dev/null || true
	@pkill -9 -f "main" 2>/dev/null || true
	@# Kill by port - try multiple methods
	@PORT=$$(if [ -f .env ]; then . ./.env && echo $${PORT:-8080}; else echo 8080; fi) && \
		echo "🔥 Killing EVERYTHING on port $$PORT..." && \
		lsof -ti:$$PORT 2>/dev/null | xargs -r kill -9 2>/dev/null || true && \
		fuser -k $$PORT/tcp 2>/dev/null || true && \
		netstat -tlnp 2>/dev/null | grep :$$PORT | awk '{print $$7}' | cut -d/ -f1 | grep -E '^[0-9]+$$' | xargs -r kill -9 2>/dev/null || true
	@# Wait and verify
	@sleep 1
	@PORT=$$(if [ -f .env ]; then . ./.env && echo $${PORT:-8080}; else echo 8080; fi) && \
		if netstat -tln 2>/dev/null | grep -q :$$PORT; then \
			echo "⚠️  Port $$PORT still in use, trying harder..." && \
			sudo fuser -k $$PORT/tcp 2>/dev/null || true && \
			sudo netstat -tlnp 2>/dev/null | grep :$$PORT | awk '{print $$7}' | cut -d/ -f1 | grep -E '^[0-9]+$$' | xargs -r sudo kill -9 2>/dev/null || true; \
		fi
	@PORT=$$(if [ -f .env ]; then . ./.env && echo $${PORT:-8080}; else echo 8080; fi) && \
		if netstat -tln 2>/dev/null | grep -q :$$PORT; then \
			echo "❌ STILL OCCUPIED! Manual intervention required."; \
			netstat -tlnp | grep :$$PORT || true; \
		else \
			echo "✅ Port $$PORT is FREE!"; \
		fi

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
