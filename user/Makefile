start:
	@bash -c "$(MAKE) -s build start-server"

build:
	@echo "  →  Building binary..."
	@go build -o user-svc main.go

start-server:
	./user-svc
