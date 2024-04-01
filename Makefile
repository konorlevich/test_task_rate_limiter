.PHONY: run stop
.SILENT:run stop

all: run

run:
	@docker compose -f deployment/docker-compose.yaml up -d --build 2>&1 1>/dev/null && docker logs -f server

recreate:
	@docker compose -f deployment/docker-compose.yaml up -d --force-recreate --build 2>&1 1>/dev/null && docker logs -f server

stop:
	@docker compose -f deployment/docker-compose.yaml down
