include .env ./api/.env
export

generate-ssl:
	@chmod +x ./scripts/setup-ssl.sh && ./scripts/setup-ssl.sh -d $(SSL_DOMAIN)

copy-certs:
	@cp -r ./certs/ ./api/certs && cp -r ./certs/ ./player/certs

start: copy-certs
	@docker compose -f docker-compose.yml up --build

start-bg: copy-certs
	@docker compose -f docker-compose.yml up --build --detach

start-dev: copy-certs
	@docker compose -f docker-compose.dev.yml up --build

start-dev-bg: copy-certs
	@docker compose -f docker-compose.dev.yml up --build --detach

destroy:
	@docker compose -f docker-compose.yml down --rmi all

destroy-dev:
	@docker compose -f docker-compose.dev.yml down --rmi all

migration-create-go:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db create_go $(name)

migration-create-sql:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db create_sql $(name)

migrate:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db migrate

rollback:
	@cd api && DB_FILE=$(DB_FILE) go run ./cmd/migrations/main.go db rollback

test:
	@cd api && go test ./...
