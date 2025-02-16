start:
	@docker compose -f docker-compose.yml up --build

destroy:
	@docker compose down --rmi all