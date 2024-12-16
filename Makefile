# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

#!make
include configs/local.env

dev:
	set -a && \
		source configs/local.env && \
	set +a && \
	nodemon --exec go run --tags dynamic $(shell pwd)/cmd/server/main.go --signal SIGTERM

run: ### Run server
	set -a && \
		source .env && \
	set +a && \
	go run --tags dynamic $(shell pwd)/cmd/server/main.go

deps: ### Install dependencies
	go mod tidy -v
	go mod vendor -v

swag-init: ### swag init
	go run cmd/docs/main.go init --parseDependency --parseInternal --parseVendor -g cmd/server/main.go

lint: ### Run Go linter
	test -s ./golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b "$(PWD)" v1.61.0
	./golangci-lint run --config .golangci.yml --timeout 10m -v --modules-download-mode vendor

migrate-up: ### Run migration up
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD_SECRET)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable" -path migrations -verbose up

migrate-down: ### Run migration down
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD_SECRET)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable" -path migrations -verbose down

new-migration: ### Generate a new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

docker:
	docker compose up -d 
	
docker-build:
	docker compose up -d --build