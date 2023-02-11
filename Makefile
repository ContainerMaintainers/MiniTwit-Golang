start:
	docker compose up

stop:
	docker compose down

fresh_start:
	docker compose up --build -d