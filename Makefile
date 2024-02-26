build:
	go build .
clean-install: clean build install-local-k8s
install-local-k8s: build
	./setup
clean:
	ctlptl delete registry ctlptl-registry || true
	minikube stop --all || true
	minikube delete --all --purge || true
	echo "Pruning docker"
	docker system prune -f --all