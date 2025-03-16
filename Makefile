include .env
export

prometheus:
	docker run -d --net=host --name=prometheus -v $(PWD)/deploy/prometheus:/etc/prometheus prom/prometheus

grafana:
	docker run -d --net=host --name=grafana -v grafana_data:/var/lib/grafana -v ${PWD}/deploy/grafana:/etc/grafana/provisioning/ -v ${PWD}/deploy/grafana/dashboards:/var/lib/grafana/dashboards -e GF_SECURITY_ADMIN_PASSWORD=admin -e GF_SECURITY_ADMIN_USER=admin grafana/grafana-oss

promdown:
	docker stop prometheus
	docker rm prometheus

grafdown:
	docker stop grafana
	docker rm grafana

run:
	source .env && go run cmd/main.go

compose-up:
	docker compose up -d 

compose-down:
	docker compose down

pg-up:
	docker compose up -d postgres-db
pg-down:
	docker compose down postgres-db

$(shell mkdir -p migrations)
DB_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres-db:5432/${POSTGRES_DB}?sslmode=disable
.PHONY: migrate-create
migrate-create:
	docker compose run --rm migrator create -ext sql -dir /migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	docker compose run --rm migrator \
	-path=/migrations \
    -database="${DB_URL}" \
	up

.PHONY: migrate-down
migrate-down:
	docker compose run --rm migrator \
        -path=/migrations \
        -database="${DB_URL}" \
        down

.PHONY: migrate-force
migrate-force:
	docker compose run --rm migrator force $(version)

.PHONY: migrate-version
migrate-version:
	docker compose run --rm migrator version