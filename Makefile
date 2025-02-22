run:
	make swag
	go build -v -o main ./cmd/user-api/main.go
	.\main

swag:
	swag init -g ./internal/http/handlers/router.go

d:
	docker compose up
jwt:
	go run cmd/jwt-token-generator/main.go

prod:
	 docker compose -f docker-compose.yml -f docker-compose.traefik.yml up --build -d