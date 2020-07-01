start:
	@bash -c "$(MAKE) -s build start-server"

build:
	@echo "  →  Building binary..."
	@go build -o auth-svc main.go

start-server:
	./auth-svc
