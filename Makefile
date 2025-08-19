ENV_FILE=infra/.env
COMPOSE=infra/docker-compose.dev.yml

up:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE) up -d --build

down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE) down

logs:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE) logs -f api

ps:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE) ps
