
build_push_order:
	docker build -f ./DockerfileOrder . --platform linux/amd64 -t guronn/otus_dz7_order:v1 --progress=plain
	docker push guronn/otus_dz7_order:v1

build_push_store:
	docker build -f ./DockerfileStore . --platform linux/amd64 -t guronn/otus_dz7_store:v1 --progress=plain
	docker push guronn/otus_dz7_store:v1

build_push_delivery:
	docker build -f ./DockerfileDelivery . --platform linux/amd64 -t guronn/otus_dz7_delivery:v1 --progress=plain
	docker push guronn/otus_dz7_delivery:v1

build_push_pay:
	docker build -f ./DockerfilePay . --platform linux/amd64 -t guronn/otus_dz7_pay:v1 --progress=plain
	docker push guronn/otus_dz7_pay:v1


build_push: build_push_order build_push_store build_push_delivery build_push_pay

#kubectl create ns 07dz
i:
	helm install 07dz .\07dz\ --namespace=07dz

ui:
	helm uninstall 07dz --namespace=07dz

ii: ui i

u:
	helm upgrade 07dz .\07dz\ --namespace=07dz

#kc get svc -n 07dz
