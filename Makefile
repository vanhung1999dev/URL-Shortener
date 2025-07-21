# Makefile for Docker Compose: build, run, stop, clean

PROJECT_NAME=URL-Shortener

.PHONY: up down clean

## Build and start the containers (no cache, detached)
up:
	@echo "🚀 Rebuilding and starting $(PROJECT_NAME)..."
	docker-compose build --no-cache
	docker-compose up -d

## Stop and remove containers (but keep volumes)
down:
	@echo "🛑 Stopping $(PROJECT_NAME)..."
	docker-compose down

## Stop and remove everything: containers + volumes
clean:
	@echo "🧹 Cleaning up all containers, networks, and volumes..."
	docker-compose down -v --remove-orphans
