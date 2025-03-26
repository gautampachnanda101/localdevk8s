#!/bin/sh
function setupDns() {
  DNS_DOMAIN=$(yq eval '.cluster.domain' config.yaml)
  CLUSTER_NAME=$(yq eval '.cluster.name' config.yaml)
  echo "Setting up domain $DNS_DOMAIN cluster profile $CLUSTER_NAME"
#  exit 1
#  minikube addons enable ingress
#  minikube addons enable ingress-dns
#  brew install chipmk/tap/docker-mac-net-connect
#  sudo brew services start chipmk/tap/docker-mac-net-connect
#  cat <<EOF | sudo tee /etc/resolver/minikube-$DNS_DOMAIN
#    domain $DNS_DOMAIN
#    nameserver $(minikube ip)
#    search_order 1
#    timeout 5
#  EOF
}
setupDns