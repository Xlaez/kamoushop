start: 
	docker compose up && go run cmd/main.go
run:
	go run cmd/main.go
init-swagger:
	swagger init -g pkg/server/server.go
.PHONY: start, run, init-swagger