.PHONY: help dev dev-build prod prod-build down down-v restart logs logs-backend logs-frontend shell-backend shell-frontend shell-db shell-redis clean run-frontend run-backend migrate-up migrate-down

# Hiển thị menu help
help:
	@echo "================================================================"
	@echo "🔥 DOCKER MAKEFILE - VIBE GO REFINE PROJECT 🔥"
	@echo "================================================================"
	@echo "🚀 Môi trường Development (Docker):"
	@echo "  make dev           - Khởi chạy dev containers (chạy ngầm)"
	@echo "  make dev-build     - Build lại images và khởi chạy dev containers"
	@echo ""
	@echo "🚀 Môi trường Production (Docker):"
	@echo "  make prod          - Khởi chạy prod containers (chạy ngầm)"
	@echo "  make prod-build    - Build lại images và khởi chạy prod containers"
	@echo ""
	@echo "💻 Chạy Local (Trên máy Host):"
	@echo "  make run-frontend  - Chạy frontend local (npm run dev)"
	@echo "  make run-backend   - Chạy backend local (go run main.go)"
	@echo "  make migrate-up    - Chạy DB migrations (cần cài golang-migrate)"
	@echo "  make migrate-down  - Rollback toàn bộ DB migrations"
	@echo ""
	@echo "🛑 Dừng & Dọn dẹp (Docker):"
	@echo "  make down          - Dừng tất cả containers"
	@echo "  make down-v        - Dừng containers VÀ XÓA VOLUMES (cảnh báo mất data)"
	@echo "  make restart       - Restart lại dev containers"
	@echo "  make clean         - Xóa toàn bộ docker cache, images rác, unused volumes"
	@echo ""
	@echo "📊 Xem Logs (Môi trường Dev):"
	@echo "  make logs          - Xem logs của tất cả services"
	@echo "  make logs-backend  - Xem logs của service backend"
	@echo "  make logs-frontend - Xem logs của service frontend"
	@echo ""
	@echo "💻 Truy cập Shell Container (Môi trường Dev):"
	@echo "  make shell-backend - Truy cập vào shell của backend"
	@echo "  make shell-frontend- Truy cập vào shell của frontend"
	@echo "  make shell-db      - Truy cập vào db (MySQL CLI)"
	@echo "  make shell-redis   - Truy cập vào redis (Redis CLI)"
	@echo "================================================================"

# --- LOCAL DEVELOPMENT (Host Machine) ---
run-frontend:
	cd apps/frontend && npm run dev

run-backend:
	cd apps/backend && go run main.go

migrate-up:
	docker exec -it vibe_backend_dev migrate -path database/migrations -database "mysql://root:root@tcp(db:3306)/vibe_db" up

migrate-down:
	docker exec -it vibe_backend_dev migrate -path database/migrations -database "mysql://root:root@tcp(db:3306)/vibe_db" down -all

# --- DOCKER DEVELOPMENT ---
dev:
	docker compose -f docker-compose.dev.yaml up -d

dev-build:
	docker compose -f docker-compose.dev.yaml up -d --build

# --- DOCKER PRODUCTION ---
prod:
	docker compose -f docker-compose.prod.yaml up -d

prod-build:
	docker compose -f docker-compose.prod.yaml up -d --build

# --- STOP & CLEAN ---
down:
	docker compose -f docker-compose.dev.yaml down
	docker compose -f docker-compose.prod.yaml down

down-v:
	docker compose -f docker-compose.dev.yaml down -v
	docker compose -f docker-compose.prod.yaml down -v

restart: down dev

clean:
	docker system prune -a --volumes -f

# --- LOGS ---
logs:
	docker compose -f docker-compose.dev.yaml logs -f

logs-backend:
	docker compose -f docker-compose.dev.yaml logs -f backend

logs-frontend:
	docker compose -f docker-compose.dev.yaml logs -f frontend

# --- SHELL ACCESS ---
shell-backend:
	docker exec -it vibe_backend_dev sh

shell-frontend:
	docker exec -it vibe_frontend_dev sh

shell-db:
	docker exec -it vibe_db_dev mysql -u root -proot vibe_db

shell-redis:
	docker exec -it vibe_redis_dev redis-cli
