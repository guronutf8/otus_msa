start:
	docker run -d -p 8000:8000 --name otus_dz2 guronn/otus_dz2

run: build
	docker run -d -p 8000:8000 --name otus_dz2 guronn/otus_dz2

build_docker:
	docker build . --platform linux/amd64 -t guronn/otus_dz2

build_nc:
	docker build --no-cache . --platform linux/amd64 -t guronn/otus_dz2

b2:
	docker build . --platform linux/amd64 -t guronn/otus_service:v6
	docker push guronn/otus_service:v6