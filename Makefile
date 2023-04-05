build-debug:
	docker build -t poker:debug -f build/package/debug/Dockerfile .
	kind load docker-image poker:debug

debug: build-debug
	helm uninstall poker
	helm install poker -f deployments/debug/values-debug.yaml deployments/poker

build:
	docker build -t poker:latest -f build/package/Dockerfile .
	kind load docker-image poker:latest

run: build
	if helm ls | grep poker ; then helm uninstall poker ; fi
	helm install poker -f deployments/poker/values.yaml --set image.repository=poker,image.tag=latest,imagegroom.repository=poker,imagegroom.tag=latest deployments/poker

clean:
	docker image rm poker:latest poker:debug
	if helm ls | grep poker ; then helm uninstall poker ; fi
