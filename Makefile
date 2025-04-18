include .env
export

generate-ssl:
	@chmod +x ./scripts/setup-ssl.sh && ./scripts/setup-ssl.sh -d $(SSL_DOMAIN)

copy-certs:
	@cp -r ./certs/ ./api/certs && cp -r ./certs/ ./player/certs

install-host-handler:
	@cd ./api && \
	CGO_ENABLED=0 GOOS=linux go build -o "${GOPATH}/bin/host_handler" ./cmd/host_handler/main.go && \
	chmod +x "${GOPATH}/bin/host_handler"
	which host_handler

build: copy-certs
	@docker compose -f docker-compose.yml build

start:
	@docker compose -f docker-compose.yml up -d

start-dev: copy-certs
	@docker compose -f docker-compose.dev.yml up --build

start-dev-bg: copy-certs
	@docker compose -f docker-compose.dev.yml up --build -d

destroy:
	@docker compose -f docker-compose.yml down --rmi all

destroy-dev:
	@docker compose -f docker-compose.dev.yml down --rmi all

migration-create-go:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db create_go $(name)

migration-create-sql:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db create_sql $(name)

migrate-init:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db init

migrate:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db migrate

rollback:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db rollback

test:
	@cd api && go test ./...
