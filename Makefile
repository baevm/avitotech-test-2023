include .env
 
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


#=============================================================================
#========================== SERVER ===========================================
#=============================================================================

## run: runs app
.PHONY: run
run:
	go run ./cmd/api

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

## build/api: build app binary for windows and linux
.PHONY: build/api
build/api:
	@echo 'Building cmd/api'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin ./cmd/api

## test: runs tests
.PHONY: test
test:
	go test -v ./... 

#=============================================================================
#========================== DATABASE =========================================
#=============================================================================

## db/migrate: executes all migrations from migrations folder
.PHONY: db/migrate
db/migrate:
	migrate -path=./internal/db/migrations -database=${DB_DSN} up

## db/migrate-down: executes all drop migrations from migrations folder
.PHONY: db/migrate-down
db/migrate-down:
	migrate -path=./internal/db/migrations -database=${DB_DSN} down

## db/unmigrate: executes one migration drop from migrations folder
.PHONY: db/migrate-down-one
db/migrate-down-one:
	migrate -path=./internal/db/migrations -database=${DB_DSN} down 1

## db/migration name=$1: creates new migration with given name
.PHONY: db/migration
db/migration:
	migrate create -seq -ext .sql -dir ./internal/db/migrations ${name}



#=============================================================================
#========================== DEPLOY ===========================================
#=============================================================================
## compose: runs docker compose
.PHONY: compose
compose:
	docker compose -f ./deploy/docker-compose.yml --env-file .env up

## docker/build: builds docker image
.PHONY: docker/build
docker/build:
	docker build --no-cache -t avito-segments-api -f ./deploy/Dockerfile .


#=============================================================================
#========================== SWAGGER ==========================================
#=============================================================================
## swagger: generates swagger docs using swaggo
.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go



#=============================================================================
#========================== OTHER ============================================
#=============================================================================
## asynq-gui: runs asynq gui to monitor queues
.PHONY: asynq-gui
asynq-gui:
	docker run -d --rm \
    --name asynqmon \
    -p 8081:8080 \
    hibiken/asynqmon --redis-addr=host.docker.internal:6379
