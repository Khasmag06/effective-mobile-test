include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d
.PHONY: compose-up

linter: ### check by golangci linter
	golangci-lint run
.PHONY: linter

migrate-create:  ### create new migration
	migrate create -ext sql -dir ./migrations -seq init_db
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	echo "y" | migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' down
.PHONY: migrate-down

test: ### run test
	go test -v ./...

cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out
.PHONY: coverage

swag: ### generate swagger docs
	swag init -g cmd/app/main.go