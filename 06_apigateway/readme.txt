minikube start --listen-address='0.0.0.0'
minikube tunnel
Set-Alias -option allscope -scope global -Name ks -Value kubectl

kubectl create ns dz6
kubectl config set-context --current --namespace=dz6
