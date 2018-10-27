build:
	docker build -t mwaaas/newrelic_prometheus_exporter:latest .
	docker push mwaaas/newrelic_prometheus_exporter:latest