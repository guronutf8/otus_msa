protoc = "c:\protoc\include\\"
annotations = "c:\Users\Noname\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway\v2@v2.16.0\protoc-gen-openapiv2\options\\"
googleapis = "C:\Users\Noname\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.16.0\third_party\googleapis\\"
schema = "f:\backup\projects\tarkovloot\protoSchema\\"
#there = "./proto"
gen:
	protoc -I . \
    	-I .. \
    	-I${protoc} \
    	-I${annotations} \
    	-I${googleapis} \
    	-I${schema} \
    	--go_out ./eshop/internal/schemagen \
    	--go-grpc_out ./eshop/internal/schemagen \
    	schema/*.proto


build_push_order:
	docker build -f ./DockerfileOrder . --platform linux/amd64 -t guronn/otus_dz9_order:v1 --progress=plain
	docker push guronn/otus_dz9_order:v1

build_push_billing:
	docker build -f ./DockerfileBilling . --platform linux/amd64 -t guronn/otus_dz9_billing:v1 --progress=plain
	docker push guronn/otus_dz9_billing:v1

build_push_notify:
	docker build -f ./DockerfileNotify . --platform linux/amd64 -t guronn/otus_dz9_notify:v1 --progress=plain
	docker push guronn/otus_dz9_notify:v1

build_push: build_push_order build_push_billing build_push_notify

#kubectl create ns 09dz
i:
	helm install 09dz .\09dz\ --namespace=09dz

ui:
	helm uninstall 09dz --namespace=09dz

ii: ui i

u:
	helm upgrade 09dz .\09dz\ --namespace=09dz