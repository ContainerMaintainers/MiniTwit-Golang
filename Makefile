start:
	docker compose up -d

stop:
	docker compose down

fresh_start:
	docker compose up --build -d