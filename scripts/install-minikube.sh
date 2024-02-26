#!/bin/sh
export CLUSTER_NAME='localdev'
export  CLUSTER_DOMAIN="test"
argument=("$@")
if [ -n "${argument[0]}" ]; then
  echo
  echo "Using cluster name: $1"
  export CLUSTER_NAME=$1
fi
if [ -n "${argument[1]}" ]; then
  echo
  echo "Using domain: $2"
  export CLUSTER_DOMAIN=$2
fi
ctlptl delete cluster $CLUSTER_NAME --cascade true || true
minikube delete --purge --all || true
ctlptl create cluster minikube --name $CLUSTER_NAME --minikube-start-flags="--driver=docker,--addons=ingress,metrics-server,dashboard,--install-addons=true"
brew install chipmk/tap/docker-mac-net-connect || true
sudo brew services start chipmk/tap/docker-mac-net-connect || true
sudo brew services
TEST_APP_URL="https://raw.githubusercontent.com/kubernetes/minikube/master/deploy/addons/ingress-dns/example/example.yaml"
FILE_CONTENTS=$(curl $TEST_APP_URL -o test.yaml)
#sed -i "s\.test\.$CLUSTER_DOMAIN\g" test.yaml
#REPLACE=".$CLUSTER_DOMAIN"
CLUSTER_DOMAIN="TT"
cp test.yaml test-back.yaml
sed "s/.test/.$CLUSTER_DOMAIN/" test.yaml | kubectl apply -f -
#cat test.yaml | tr '.test' $REPLACE | echo
#echo $FILE_CONTENTS
#echo $FILE_CONTENTS >  test.yaml

cat test.yaml | kubectl apply -f -
#echo $TEST_APP | kubectl apply -f -
#kubectl apply -f https://raw.githubusercontent.com/kubernetes/minikube/master/deploy/addons/ingress-dns/example/example.yaml
nslookup hello-john.$CLUSTER_DOMAIN $(minikube ip) || true

# Add minikube to host machine dns lookup
sudo rm -rf /etc/resolver/minikube-$CLUSTER_DOMAIN || true
sudo touch /etc/resolver/minikube-$CLUSTER_DOMAIN
USAGE="domain $CLUSTER_DOMAIN
nameserver $(minikube ip)
search_order 11
timeout 5"

echo $USAGE | sudo tee -a /etc/resolver/minikube-$CLUSTER_DOMAIN
#read -r -d '' RESOLVER_CONF << EOM
#domain $CLUSTER_DOMAIN
#nameserver $(minikube ip)
#search_order 11
#timeout 5
#EOM
#echo $RESOLVER_CONF
#echo $RESOLVER_CONF > /etc/resolver/minikube-$CLUSTER_DOMAIN
sudo cat /etc/resolver/minikube-$CLUSTER_DOMAIN

# Test to actually reach the ingress without specifying dns
curl hello-john.$CLUSTER_DOMAIN