
.PHONY=build
build:
	docker build -t redis-custom-exporter:1.0.0 .

load:
	kind load docker-image --name redis redis-custom-exporter:1.0.0

deploy: build load
	kubectl delete -f manifests/ || true && kubectl apply -f manifests/
