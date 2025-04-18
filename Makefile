include .env

.PHONY: debug dependency protoc-gen clean-proto migration migrate-up migrate-down run-test

debug:
	@echo "Starting the engine..."
	go run "${SERVICE_DIR}" service --rpc "${RPC_PORT}" --gateway "${GATEWAY_PORT}"

dependency:
	@echo "Downloading all Go dependencies..."
	go mod download
	go mod verify
	go mod tidy
	@echo "All Go dependencies are downloaded. You can now run 'make debug' to compile locally or 'docker compose up -d' to build the app."

proto-gen: clean-proto
	@echo "Generating the gRPC stubs..."
	./scripts/protoc-gen.sh
	@echo "Stubs generated successfully. All stubs are in the 'stubs/' directory."
	@echo "Generating Swagger UI..."
	./scripts/swagger-ui-gen.sh
	@echo "Swagger UI generated successfully. For older versions, check './cache/swagger-ui' directory."
	@echo "DO NOT EDIT ANY FILES IN THE STUBS MANUALLY!"

clean-proto:
	@echo "Removing all previously generated stubs..."
	rm -rf stubs/article/*
	@echo "All stubs have been cleaned up successfully."

migration:
	@echo "Creating a new named migration..."
	migrate create -ext sql -dir db/migrations "$(name)"
	@echo "Migration file created successfully."

migrate-up:
	@echo "Applying migrations (up)..."
	migrate -database "${MIGRATE_DNS}" -path db/migrations up
	@echo "Migrations applied successfully."

migrate-down:
	@echo "Rolling back migrations (down)..."
	migrate -database "${MIGRATE_DNS}" -path db/migrations down
	@echo "Migrations rolled back successfully."

run-test:
	@echo "Running all tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo "All tests executed successfully."
