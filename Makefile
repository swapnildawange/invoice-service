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

run : 
	go run ./cmd/main.go

createdb:
	docker exec -it invoice_service_postgres_1 createdb --username=postgres --owner=postgres invoicing

opendb:
	docker exec -it invoice_service_postgres_1 psql -U postgres invoicing

dropdb:
	docker exec -it invoice_service_postgres_1 dropdb  --username=postgres invoicing


migrateup:
	migrate -path db/migration -database "postgres://postgres:password@localhost/invoicing?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:password@localhost/invoicing?sslmode=disable" -verbose down

# createmigration:
# 	migrate create -ext sql -dir db/migration -seq init_schema