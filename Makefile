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
	go run cmd/main.go

composeup:
	docker compose up -d 

composedown:
	docker compose down