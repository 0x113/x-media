start:
	@bash -c "$(MAKE) -s build swagger start-server"

build:
	@echo "  →  Building binary..."
	@go build -o tvshow main.go

start-server:
	./tvshow

swagger:
	@echo "  →  Generating docs..."
	@swag i
