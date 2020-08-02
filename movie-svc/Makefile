start:
	@bash -c "$(MAKE) -s build swagger start-server"

build:
	@echo "  →  Building binary..."
	@go build -o movie-svc main.go

start-server:
	./movie-svc

swagger:
	@echo "  →  Generating docs..."
	@swag i
