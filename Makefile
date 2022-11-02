start: 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"
	
stop:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

run : start
	go run ./cmd/main.go