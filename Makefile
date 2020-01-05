include .env_make

start:
	@bash -c "$(MAKE) -s build start-server"

build:
	@echo "  >  Building binary..."
	@go build

start-server:
	@echo "  >  Starting proejct at localhost:$(PORT)..."
	@DB_USER=$(DB_USER) DB_PASS=$(DB_PASS) DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_NAME=$(DB_NAME) JWT_SECRET=$(JWT_SECRET) VIDEO_DIR=$(VIDEO_DIR) MOVIES_SUB_DIR=$(MOVIES_SUB_DIR) FRONTEND_DIR=$(FRONTEND_DIR) PORT=$(PORT) ./x-media
