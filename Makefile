version = latest
deploy:
	docker build -t mwaaas/newrelic_prometheus_exporter:$(version) .
	docker push mwaaas/newrelic_prometheus_exporter:$(version)

deploy_target:
	docker build --target build-env -t mwaaas/newrelic_prometheus_exporter:latest_build_env .
	docker push mwaaas/newrelic_prometheus_exporter:latest_build_env

deploy_all: deploy deploy_target