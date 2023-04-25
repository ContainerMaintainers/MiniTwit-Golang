include .env

start:
	docker compose up

stop:
	docker compose down

fresh_start:
	docker compose up --build

prune_docker:
	docker system prune -a