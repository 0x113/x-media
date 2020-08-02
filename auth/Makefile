start:
	@bash -c "$(MAKE) -s build swagger start-server"

build:
	@echo "  →  Building binary..."
	@go build -o auth-svc main.go

start-server:
	./auth-svc

swagger:
	@echo "  →  Generating docs..."
	@swag i
