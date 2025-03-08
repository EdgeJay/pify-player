include .env
export

generate-ssl:
	@chmod +x ./scripts/setup-ssl.sh && ./scripts/setup-ssl.sh -d $(SSL_DOMAIN)

start:
	@docker compose -f docker-compose.yml up --build

start-dev:
	@docker compose -f docker-compose.dev.yml up --build

start-dev-bg:
	@docker compose -f docker-compose.dev.yml up --detach

destroy:
	@docker compose -f docker-compose.yml down --rmi all

destroy-dev:
	@docker compose -f docker-compose.dev.yml down --rmi all
