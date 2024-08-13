.PHONY: docker k8s etcd

docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f  orangeqwj/webook:v0.0.1 
	@docker build -t orangeqwj/webook:v0.0.1 .

k8s:
	@kubectl delete deployment webook
	@kubectl apply -f ./k8s-webook-deployment.yaml
	@kubectl apply -f ./k8s-webook-service.yaml 
	@kubectl apply -f ./k8s-mysql-deployment.yaml
	@kubectl apply -f ./k8s-mysql-service.yaml
	@kubectl apply -f ./k8s-mysql-pvc.yaml
	@kubectl apply -f ./k8s-mysql-pv.yaml
	@kubectl apply -f ./k8s-redis-deployment.yaml
	@kubectl apply -f ./k8s-redis-service.yaml
	@kubectl apply -f ./k8s-ingress-nginx.yaml
etcd:
	etcdctl --endpoints=127.0.0.1:12379 put /webook "$$(cat config/dev.yaml)"


run: docker k8s

