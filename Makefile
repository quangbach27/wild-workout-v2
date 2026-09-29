.PHONY: up down down-clean

up:
	docker compose up --build -d
	docker compose logs -f wild-workout-backend

down:
	docker compose down

# Also removes volumes (wipes Postgres data) and orphans
down-clean:
	docker compose down --volumes --remove-orphans
