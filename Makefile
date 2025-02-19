include .env
export

generate-ssl:
	@chmod +x ./scripts/setup-ssl.sh && ./scripts/setup-ssl.sh -d $(DOMAIN)

start:
	@docker compose -f docker-compose.yml up --build

start-dev:
	@docker compose -f docker-compose.dev.yml up --build

destroy:
	@docker compose -f docker-compose.yml down --rmi all

destroy-dev:
	@docker compose -f docker-compose.dev.yml down --rmi all
