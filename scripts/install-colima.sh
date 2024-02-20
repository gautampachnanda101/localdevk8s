#!/bin/sh
export CLUSTER_NAME='localdev'
function deleteColima() {
    colima list | awk '{print $1}' | grep -v 'PROFILE' | xargs colima stop --force --profile || true
    colima list | awk '{print $1}' | grep -v 'PROFILE' | xargs colima delete --force --profile || true
    rm -rf ~/.lima ~/.colima || true
}
argument=("$@")
if [ -n "${argument[0]}" ]; then
  echo
  echo "Using cluster name: $1"
  export CLUSTER_NAME=$1
fi
deleteColima
START_CMD="colima start --profile $CLUSTER_NAME --kubernetes --cpu 4 --memory 8 --disk 100 --dns 1.1.1.1 --dns 8.8.8.8 --activate"
echo "Using command $START_CMD"
$START_CMD