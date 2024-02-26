# localdevk8s
A simple minikube cluster for local dev

## Prerequisites
 - Docker Desktop
 - Mac
 - Homebrew
 - Golang
### Install Prerequisites

Install [homebrew](https://brew.sh/)

```agsl
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
(echo; echo 'eval "$(/opt/homebrew/bin/brew shellenv)"') >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
brew update&& brew install golang
```
[GVM](https://github.com/moovweb/gvm) for managing go versions

```agsl
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```

Install [docker desktop](https://docs.docker.com/desktop/install/mac-install/)

Configure docker desktop to ensure rosetta emulator is enabled. Sufficient resources are given to match the cluster [config](./config.yaml) for [minikube](https://minikube.sigs.k8s.io/docs/start/)

## Starting the cluster
Run the make command
```agsl
make install-local-k8s
```
This will also install some additional tools like 
    - kubectx
    - openlens
    - yq
    - jq
    - stern
    - minikube
    - kubectl
    - docker
    - colima
    - cmctl
    - mkcert
    - k9s
    - helm
    - tilt
    - ctlptl

**You will be prompted for password a few times to configure certs and cluster**

To destroy the cluster run command
```agsl
make clean
```

### Troubleshooting
A restart of machine would need a cluster rebuild i.e.
Run
```agsl
make clean-install
```
This can take a few minutes as we are spinning up a new cluster

**You will be prompted for password a few times to configure certs and cluster**

once cluster is running
You can access
| Tool    | UTL | NOTES |
| -------- | ------- | ------- |
| Minikube Dashboard  | https://kubernetes-dashboard.localhost    |  replace .localhost with clusterDomain |
| Minikube Dashboard  | https://traefik.localhost/dashboard/    |  replace .localhost with clusterDomain |
| Test App    | https://hello-john.localhost   | replace .localhost with clusterDomain |

### Configuration
```agsl
k8Provider: minikube
clusterDomain: gautamlocalhost
clusterName: gautam
clusterNodes: 1
```

### Caveats
The currently implementation is only functional on minikube and mac os.