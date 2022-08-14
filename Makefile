TAG=$(shell git describe --tags)
CLUSTER=kind-scout

run:
	@ go run cmd/main.go

lint:
	@ golangci-lint run

up:
	@ kubectl cluster-info --context kind-$(CLUSTER)
	@ docker build . -t scout:dev -f dockerfiles/scout/Dockerfile
	@ docker build . -t scout-portal:dev -f dockerfiles/portal/Dockerfile
	@ kind load docker-image scout:dev scout-portal:dev --name $(CLUSTER)
	@ kubectl apply -f ./k8s/deployment-dev.yaml

down:
	@ kubectl delete -f ./k8s/deployment-dev.yaml

portal:
	@ kubectl exec -n scout -it portal -- bash

# https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#create-the-cluster
setup-kind:
	@ kind create cluster --name $(CLUSTER) --config=./cluster/kind-no-cni-config.yaml
	@ cilium install

teardown-kind:
	@ kind delete cluster --name $(CLUSTER)

# https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/#install-the-cilium-cli
cilium-cli:
	@ curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
	@ sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
	@ sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
	@ rm cilium-linux-amd64.tar.gz{,.sha256sum}
