PROJECT_NAME=nyanbot
VERSION?=0.0.10

testall:
	go test ./...

release:
	git tag -a v$(VERSION) -m "new tag by make"
	git push origin tag v$(VERSION)

dockerbuild:
	cd docker/ && \
	docker buildx build \
		-t commojun/nyanb:$(VERSION) \
		--platform linux/amd64,linux/arm/v7,linux/arm/v6 \
		--build-arg VERSION=$(VERSION) \
		--push \
		-f ./Dockerfile .

pods:
	kubectl get pods

init:
	kubectl apply \
	-f ./kube/namespace.yml \
	-f ./kube/redis.yml \
	-f ./kube/server.yml

deploy:
	kubectl apply \
	-f ./kube/server.yml \
	-l deploy

secret:
	-kubectl delete secret nyan-secret
	kubectl create secret generic \
		--save-config nyan-secret \
		--from-env-file ./envfile

redis-cli:
	kubectl exec -it redis redis-cli

hello:
	kubectl apply -f kube/hello.yml

export:
	-kubectl delete -f kube/export.yml
	kubectl apply -f kube/export.yml
