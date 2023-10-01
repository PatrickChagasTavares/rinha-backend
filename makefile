.PHONY: setup build docker-up docker-down up-local down-local run-gin run-echo docs mocks migration-create migration-up migration-down

setup:
	@echo "installing swaggo..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "installing golang-migrate..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "installing mockgen..."
	@go install github.com/golang/mock/mockgen@latest
	@echo "downloading project dependencies..."
	@go mod download

build: 
	@echo $(NAME_IMAGE): Compilando o micro-serviço
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -o dist/$(NAME_IMAGE)/main cmd/main.go 
	@echo $(NAME_IMAGE): Construindo a imagem
	@docker build -t $(NAME_IMAGE) .

docker-up:
	@docker compose -f "docker/docker-compose.yml" up -d --build

docker-down:
	@docker compose -f "docker/docker-compose.yml" down

up-local:
	@docker compose -f "docker/db/docker-compose.yml" up -d --build
	@docker compose -f "docker/observability/docker-compose.yml" up -d --build

down-local:
	@docker compose -f "docker/db/docker-compose.yml" down
	@docker compose -f "docker/observability/docker-compose.yml" down

run-gin:
	env=local go run cmd/rinha-gin/main.go

run-echo:
	env=local go run cmd/rinha-echo/main.go

docs:
	@swag init --parseDependency -g cmd/rinha-gin/main.go

mocks:
	@go generate ./... 

migration-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

migration-up:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose up

migration-down:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose down 1