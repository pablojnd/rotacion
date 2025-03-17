.PHONY: build run docker-build docker-run clean

# Variables
APP_NAME=rotacion-api
DOCKER_IMAGE=rotacion-api:latest

# Compilación local
build:
	go build -o $(APP_NAME) .

# Ejecución local
run: build
	./$(APP_NAME)

# Construcción de Docker
docker-build:
	docker-compose build

# Ejecución de Docker
docker-run:
	docker-compose up -d

# Limpieza
clean:
	rm -f $(APP_NAME)
	go clean

# Detener contenedores Docker
docker-stop:
	docker-compose down

# Mostrar logs de Docker
docker-logs:
	docker-compose logs -f
