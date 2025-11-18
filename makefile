# ------------------------------------
# OMHS Backend - Makefile
# ------------------------------------

COMPOSE_DEV = docker compose -f docker-compose.dev.yml
COMPOSE_PROD = docker compose -f docker-compose.yml

# ----------🧩 DEVELOPMENT ----------
omhs-dev:
	$(COMPOSE_DEV) up

omhs-dev-build:
	$(COMPOSE_DEV) up --build

omhs-dev-down:
	$(COMPOSE_DEV) down

omhs-dev-logs:
	$(COMPOSE_DEV) logs -f backend

omhs-dev-restart:
	$(COMPOSE_DEV) down
	$(COMPOSE_DEV) up

# ----------🚀 PRODUCTION ----------
omhs-prod:
	$(COMPOSE_PROD) up -d

omhs-prod-build:
	$(COMPOSE_PROD) up --build -d

omhs-prod-down:
	$(COMPOSE_PROD) down

omhs-prod-logs:
	$(COMPOSE_PROD) logs -f backend

# ----------💻 GO COMMANDS ----------
omhs-run:
	go run main.go

omhs-build:
	go build -o omhs .

omhs-fmt:
	go fmt ./...

omhs-test:
	go test ./...

omhs-vet:
	go vet ./...

omhs-tidy:
	go mod tidy

omhs-clean:
	rm -f omhs

# ----------🧰 UTILITIES ----------
omhs-db-shell:
	docker exec -it mongo mongosh

omhs-backend-bash:
	docker exec -it backend sh

omhs-stop-all:
	docker compose down

# ----------📚 HELP ----------
omhs-help:
	@echo ""
	@echo "OMHS Commands:"
	@echo ""
	@echo "  omhs-dev            - Start dev environment"
	@echo "  omhs-dev-build      - Start dev with rebuild"
	@echo "  omhs-dev-down       - Stop dev environment"
	@echo "  omhs-dev-logs       - Follow backend logs (dev)"
	@echo "  omhs-dev-restart    - Restart dev environment"
	@echo ""
	@echo "  omhs-prod           - Start production (detached)"
	@echo "  omhs-prod-build     - Production build & run"
	@echo "  omhs-prod-down      - Stop production"
	@echo "  omhs-prod-logs      - Follow backend logs (prod)"
	@echo ""
	@echo "  omhs-run            - Run Go app locally"
	@echo "  omhs-build          - Build Go app"
	@echo "  omhs-fmt            - Format Go code"
	@echo "  omhs-test           - Run tests"
	@echo "  omhs-vet            - Run go vet"
	@echo "  omhs-tidy           - Clean go.mod"
	@echo "  omhs-clean          - Remove binary"
	@echo ""
	@echo "  omhs-db-shell       - Enter Mongo shell"
	@echo "  omhs-backend-bash   - Enter backend container shell"
	@echo "  omhs-stop-all       - Stop ALL compose environments"
	@echo ""
