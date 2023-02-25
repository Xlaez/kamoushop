start: 
	docker compose up && go run cmd/main.go
run:
	go run cmd/main.go
init-swagger:
	swag init -g cmd/main.go
.PHONY: start, run, init-swagger