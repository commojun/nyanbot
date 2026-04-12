PROJECT_NAME=nyanbot
VERSION?=0.0.11

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

get-all:
	kubectl get all

pods:
	kubectl get pods -o wide --watch

init:
	kubectl apply \
	-f ./kube/namespace.yml \
	-f ./kube/server.yml \
	-f ./kube/anniversary.yml

deploy:
	kubectl apply \
	-f ./kube/server.yml \
	-l deploy

secret:
	-kubectl delete secret nyan-secret
	kubectl create secret generic \
		--save-config nyan-secret \
		--from-env-file ./envfile

hello:
	-kubectl delete -f kube/hello.yml
	kubectl apply -f kube/hello.yml

logs/%:
	kubectl logs --timestamps=true --prefix=true -f -l app=$*

logs-all:
	kubectl logs --timestamps=true --prefix=true -f -l app

shell/%:
	kubectl exec -it $* -- bash

restart/server:
	kubectl rollout restart deployment/server-deployment

# CronJobs do not need to be restarted. The next run will use the updated data.
# To run immediately, use 'kubectl create job --from=cronjob/<name> <job-name>'
restart/alarm:
	@echo "Alarm is a CronJob, no restart needed."

restart/anniversary:
	@echo "Anniversary is a CronJob, no restart needed."

restart/all:
	make restart/server
