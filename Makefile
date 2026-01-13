.PHONY: run test clean

# Run the service directly using go run
run:
	@echo "Syncing dependencies and starting service on port 8080..."
	@go mod tidy
	@go run main.go

# Send a sample test request
test:
	@echo "Testing /event-forecast..."
	@curl -s -X POST http://localhost:8080/event-forecast \
		-H "Content-Type: application/json" \
		-d '{"name":"Match","location":{"latitude":19.076,"longitude":72.877},"start_time":"2026-01-14T17:00:00","end_time":"2026-01-14T20:00:00"}' | json_pp 2>/dev/null

clean:
	@go clean
	@echo "Cleaned build artifacts."
