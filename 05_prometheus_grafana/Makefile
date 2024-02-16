start:
	docker run -d -p 8000:8000 --name otus_dz4 guronn/otus_dz4

run: build
	docker run -d -p 8000:8000 --name otus_dz4 guronn/otus_dz4

build:
	docker build . --platform linux/amd64 -t guronn/otus_dz4

build_nc:
	docker build --no-cache . --platform linux/amd64 -t guronn/otus_dz4

build_push:
	docker build . --platform linux/amd64 -t guronn/otus_dz5:v4 --progress=plain
	docker push guronn/otus_dz5:v4

lp:
	kubectl logs --selector=app=prometheus

i:
	helm install dz5 .\chart_dz5\

ui:
	helm uninstall dz5

ii: ui i