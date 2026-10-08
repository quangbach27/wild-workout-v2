.PHONY: up down down-clean seed web storybook

up:
	docker compose up --build -d
	docker compose logs -f wild-workout-backend

down:
	docker compose down

# Also removes volumes (wipes Postgres data) and orphans
down-clean:
	docker compose down --volumes --remove-orphans

# Load sample trainers and open hours (needs `make up` running)
seed:
	docker compose exec -T wild-workout-db psql -U user -d wild-workout < backend/seed/dev.sql

# Vite dev server on http://localhost:3000
web:
	cd web && npm run dev

# Storybook on http://localhost:6006
storybook:
	cd web && npm run storybook
