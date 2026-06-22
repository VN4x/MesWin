.PHONY: dev-backend dev-frontend build-backend build-frontend build-docker

dev-backend:
	cd backend && go run . serve --http=127.0.0.1:8090

dev-frontend:
	cd frontend && npm run dev

build-backend:
	cd backend && CGO_ENABLED=0 go build -o ../dist/meswin .

build-frontend:
	cd frontend && npm run build
	mkdir -p backend/pb_public
	cp -r frontend/build/* backend/pb_public/

build-docker:
	docker compose -f deploy/docker-compose.yml build
