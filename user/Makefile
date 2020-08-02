start:
	@bash -c "$(MAKE) -s build swagger start-server"

build:
	@echo "  →  Building binary..."
	@go build -o user-svc main.go

start-server:
	./user-svc

swagger:
	@echo "  →  Generating docs..."
	@swag i
