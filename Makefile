# Run docker compose up -d with this Makefile

PHONY: up down client go-run ts-run

up:
	docker compose up -d

down:
	docker compose down

client:
	docker exec -it clickhouse clickhouse-client

go-run-storer:
	cd go-app && go mod tidy && go run cmd/storer/main.go

go-run-asker:
	cd go-app && go mod tidy && go run cmd/asker/main.go

ts-run-storer:
	cd ts-app && npm install && npm run storer

ts-run-asker:
	cd ts-app && npm install && npm run asker
