start:
	@bash -c "$(MAKE) -s build start-server"

build:
	@echo "  →  Building binary..."
	@go build -o movie-svc main.go

start-server:
	./movie-svc
