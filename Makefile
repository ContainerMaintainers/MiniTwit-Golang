start:
	docker compose --env-file ./.env -f docker/docker-compose.yml up 

stop:
	docker compose --env-file ./.env -f docker/docker-compose.yml down

fresh_start:
	docker compose --env-file ./.env -f docker/docker-compose.yml up --build

prune_docker:
	docker system prune -a