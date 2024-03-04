
build_nc:
	docker build --no-cache . --platform linux/amd64 -t guronn/otus_dz6_auth

build_auth:
	docker build -f ./DockerfileAuth . --platform linux/amd64 -t guronn/otus_dz6_auth:v1 --progress=plain

build_push_auth:
	docker build -f ./DockerfileAuth . --platform linux/amd64 -t guronn/otus_dz6_auth:v2 --progress=plain
	docker push guronn/otus_dz6_auth:v2

build_push_app:
	docker build -f ./DockerfileApp . --platform linux/amd64 -t guronn/otus_dz6_app:v2 --progress=plain
	docker push guronn/otus_dz6_app:v2


i:
	helm install dz6 .\dz6\ --namespace=dz6

ui:
	helm uninstall dz6

ii: ui i

u:
	helm upgrade dz6 .\dz6\