kind:
	kind create cluster --config deployments/kind/kind-config.yaml
	helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
	helm repo update
	helm upgrade --install --set args={--kubelet-insecure-tls} metrics-server metrics-server/metrics-server --namespace kube-system
	kubectl apply -f deployments/kind/deploy-nginx4kind.yaml 

kind-delete:
	kind delete cluster

build-debug:
	docker build -t poker:debug -f build/package/debug/Dockerfile .
	kind load docker-image poker:debug

debug: build-debug
	helm uninstall poker -n poker
	helm install poker -f deployments/debug/values-debug.yaml deployments/poker -n poker

build:
	docker build -t poker:latest -f build/package/Dockerfile .
	kind load docker-image poker:latest

run: build
	kubectl create namespace poker 2> /dev/null || true
	if helm ls | grep poker ; then helm uninstall poker -n poker ; fi
	helm install poker -n poker -f deployments/poker/values.yaml --set image.repository=poker,image.tag=latest,imagegroom.repository=poker,imagegroom.tag=latest deployments/poker

clean:
	docker image rm poker:latest poker:debug 2> /dev/null || true
	if helm ls | grep poker ; then helm uninstall poker -n poker ; fi
	kubectl delete namespace poker 2> /dev/null || true
