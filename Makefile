build:
	docker build -t mwaaas/newrelic_prometheus_exporter:latest .
	docker push mwaaas/newrelic_prometheus_exporter:latest

build_target:
	docker build --target build-env -t mwaaas/newrelic_prometheus_exporter:latest_build_env .
	docker push mwaaas/newrelic_prometheus_exporter:latest_build_env