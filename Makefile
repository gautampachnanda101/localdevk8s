build:
	go build .
clean:
	ctlptl delete registry ctlptl-registry || true
	minikube stop --all || true
	minikube delete --all --purge || true
	echo "Pruning docker"
	docker system prune -f --all
install-local-k8s: build
	./setup
install: clean build install-local-k8s