start: 
	docker compose up && go run cmd/main.go
run:
	go run cmd/main.go

.PHONY: start, run