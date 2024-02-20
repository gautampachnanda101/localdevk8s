package main

import (
	"os/exec"
	 log "github.com/sirupsen/logrus"
	"strings"
	"github.com/gookit/color"
	"errors"
	"os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "bufio"
    parser "org/pachnanda/machine/setup/parser"
    "fmt"
)
var k8Provider = "Colima"
var clusterName = "localdev"
var clusterDomain = "localhost"
var installScript = ""
var minikubeInstall = ""
var minikubeTest = ""
var colimaInstall = ""
var colimaTest = ""
var clean = ""
var helmInstall = ""
var certsSetup = ""
var clusterNodes = "0"
var k8sVersion="v1.25.0"
type Config struct {
	Cluster struct {
	    Name string `yaml:"name"`
	    Nodes string `yaml:"nodes"`
	    Domain string `yaml:"domain"`
	}`yaml:"cluster"`
	Clean string `yaml:"clean"`
    Minikube struct {
        InstallCluster string `yaml:"installCluster"`
        TestCluster string `yaml:"testCluster"`
    } `yaml:"minikube"`
   Colima  struct {
         InstallCluster string `yaml:"installCluster"`
         TestCluster string `yaml:"testCluster"`
   } `yaml:"colima"`
   HelmInstall string `yaml:"helmInstall"`
   CertsSetup  string `yaml:"certsSetup"`
   K8s struct {
    Version  string `yaml:"version"`
    Provider string `yaml:"provider"`
   } `yaml:"k8s"`
}

func init() {
    log.SetFormatter(&log.JSONFormatter{})
}

func readConf(filename string) (Config) {
    yamlFile, err := ioutil.ReadFile(filename)

    if err != nil {
        panic(err)
    }

    var config Config

    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
       panic(err)
    }
    log.Println("Loading", filename)
    color.Info.Tips("K8Provider : %s\n", config.K8s.Provider)
    color.Info.Tips("ClusterDomain : %s\n", config.Cluster.Domain)
    color.Info.Tips("ClusterName : %s\n", config.Cluster.Name)
    log.Println("MinikubeInstall: ", config.Minikube.InstallCluster)
    log.Println("MinikubeTest: ", config.Minikube.TestCluster)
    log.Println("ColimaInstall: ", config.Colima.InstallCluster)
    log.Println("ColimaTest: ", config.Colima.TestCluster)
    log.Println("Clean: ", config.Clean)
    log.Println("helmInstall: ", config.HelmInstall)
    log.Println("certsSetup: ", config.CertsSetup)
    log.Println("clusterNodes: ", config.Cluster.Nodes)
    log.Println("K8s version: ", config.K8s.Version)
    return config
}
func configureTemplates(domain string) {
    type ResolverConfig struct {
        Domain string
        Ip string
    }
    ip, err := exec.Command("/bin/sh", "-c", "minikube ip -p "+clusterName).Output()
        if err != nil {
            log.Fatal(err)
        }
    minikubeIp := strings.Replace(string(ip), "\n", "", -1)
    color.Notice.Println("The ip is [%s]", minikubeIp)
    m := make(map[string]string)
    m["Domain"] = domain
    m["Ip"] = minikubeIp
    log.Println(m)
    buildResolver(m)
    buildTraefikCerts(m)
    buildTraefikValues(m)
    buildDashboardIngress(m)
    buildTestApp(m)
    buildClusterIssuer(m)

}
func startTunnel(){
    color.Info.Tips("String minikube tunnel for profile",clusterName)
    cmd := exec.Command("/bin/sh", "-c",fmt.Sprintf("minikube -p %s tunnel &",clusterName))
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
}
func startK8sDashboard(){
    cmd := exec.Command("/bin/sh", "-c",fmt.Sprintf("minikube -p %s dashboard &",clusterName))
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
}
func main() {
    c := readConf("config.yaml")
    k8Provider=c.K8s.Provider
    clusterName=c.Cluster.Name
    clusterDomain=c.Cluster.Domain
    clean=c.Clean
    helmInstall=c.HelmInstall
    certsSetup=c.CertsSetup
    clusterNodes=c.Cluster.Nodes
    k8sVersion=c.K8s.Version
    color.Info.Tips("%s install started. name: %s & domain: %s",k8Provider,clusterDomain,clusterName)
    prepDependencies()
    runShellScripts(clean,"Destroy existing "+ k8Provider)
    if k8Provider == "colima" {
        colimaInstall=c.Colima.InstallCluster
        colimaTest=c.Colima.TestCluster
        runShellScripts(colimaInstall,"Install Cluster "+ k8Provider)
        configureTemplates(clusterDomain)
        runShellScripts(certsSetup, "Setup cluster certs")
        runShellScripts(helmInstall, "Install Helm Components")

        runShellScripts(colimaTest, "Verify Cluster")
        // cleanCreateColima()
        applyTraefikCerts()
        applyTraefik()
    } else if k8Provider == "minikube" {
        minikubeInstall=c.Minikube.InstallCluster
        minikubeTest=c.Minikube.TestCluster
        runShellScripts(minikubeInstall,fmt.Sprintf("Install Cluster %s", k8Provider))
        configureTemplates(clusterDomain)
       // runBackgroundShellScript(fmt.Sprintf("minikube -p %s tunnel",clusterName),"Starting tunnel")
        runShellScripts(certsSetup, "Setup cluster certs")
        startTunnel()
        runShellScripts(helmInstall, "Install Helm Components")
//         configureTemplates(clusterDomain)
        // startTunnel()
        applyTraefikCerts()
        applyTraefik()
        applyDashboard()
        runShellScripts(minikubeTest, "Verify Cluster")
        //runBackgroundShellScript(fmt.Sprintf("minikube -p %s dashboard",clusterName), "start dashboard")

        // cleanCreateMinikube()
        // installDockerForMac()
    }else {
        color.Error.Println("K8Provider %s not supported", k8Provider)
    }

     //installCerts()
 	//installCerts()
  	//applyNginxIngress()
 	//applyTraefik()
    //  	setupHosts()
}
func executeCommand(cmdToExecute string) {
    color.Error.Println("Running background shell script", cmdToExecute)
    cmd := exec.Command("/bin/sh", "-c",cmdToExecute)
    cmd.Stdout = os.Stdout
    err := cmd.Start()
    if err != nil {
       color.Error.Println("Failed to execute " + cmdToExecute + " cause:  ", err)
    }
    color.Info.Tips(fmt.Sprintf("Just ran %s  subprocess %d, exiting\n",cmdToExecute, cmd.Process.Pid))
}

func runBackgroundShellScript(scripts string,label string){
    scriptToExecute:= scripts + " &"
    executeCommand(scriptToExecute)
}
func runShellScripts(scripts string,label string){
    color.Info.Tips("Executing ... %s" , label)
    color.Info.Tips(scripts)
    scanner := bufio.NewScanner(strings.NewReader(scripts))
    for scanner.Scan() {
        cmdToExecute := scanner.Text()
        cmdToExecute = strings.Replace(cmdToExecute,"$CLUSTER_DOMAIN",clusterDomain,-1)
        cmdToExecute = strings.Replace(cmdToExecute,"$CLUSTER_NAME",clusterName,-1)
        cmdToExecute = strings.Replace(cmdToExecute,"$CLUSTER_NODES",clusterNodes,-1)
        cmdToExecute = strings.Replace(cmdToExecute,"$K8S_VERSION",k8sVersion,-1)
        color.Info.Tips("Executing command ... %s",  cmdToExecute)
        cmd := exec.Command("/bin/sh", "-c", cmdToExecute)
        out, err := cmd.Output()
        if err != nil {
            color.Error.Println("Failed to execute [" + cmdToExecute + "] - cause: " + fmt.Sprint(err) + " output: " + string(out))
            //color.Error.Println("Failed to execute " + cmdToExecute + " cause:  ", err)
        }
        color.Info.Tips("Command output %s" , string(out))
    }

    if err := scanner.Err(); err != nil {
        color.Error.Printf("Error occurred: %v\n", err)
    }
     color.Info.Tips("Executing %s Completed", label)
}

func installDockerForMac() {
    command := []string{
        "./scripts/install-docker-net-connect.sh",
        clusterName,
    }

    execute("./scripts/install-docker-net-connect.sh", command)
}
func cleanCreateMinikube(){
    color.Info.Tips("Cleaning ... %s", k8Provider)
    command := []string{
            "./scripts/install-minikube.sh",
            clusterName,
            clusterDomain,
    }

    execute("./scripts/install-minikube.sh", command)
//     listExistingMinikube()
//     deleteExistingMinikube()
//     installMinikube()

}
func cleanCreateColima(){
    color.Info.Tips("Cleaning %s", k8Provider)
    listExistingColima()
    stopExistingColima()
    installColima()
    listExistingColima()
}
func execute(script string, command []string) (bool, error) {
    cmd := &exec.Cmd{
        Path:   script,
        Args:   command,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    }

    color.Info.Tips("Executing command ", cmd)

    err := cmd.Start()
    if err != nil {
        return false, err
    }

    err = cmd.Wait()
    if err != nil {
        return false, err
    }

    return true, nil
}
func installMinikube(){
    color.Info.Tips("Installing %s",  k8Provider)

    command := []string{
        "./scripts/install-minikube.sh",
        clusterName,
    }

    execute("./scripts/install-minikube.sh", command)
}
func installColima() {
    color.Info.Tips("Installing %s", k8Provider)

    command := []string{
        "./scripts/install-colima.sh",
        clusterName,
    }

    execute("./scripts/install-colima.sh", command)
}

func deleteExistingMinikube() {
    cmd := exec.Command("/bin/sh", "-c", "minikube delete --purge --all || true")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to delete " + k8Provider + " cause: ", err)
    }
    color.Notice.Println(string(out))
}
func stopExistingColima() {
    cmd := exec.Command( "/bin/sh", "-c", "colima list | awk '{print $1}' | grep -v 'PROFILE' | xargs colima stop --force -p || true")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to stop " + k8Provider + " cause: ", err)
    }

    color.Notice.Println(string(out))

    cmd = exec.Command( "/bin/sh", "-c", "colima list | awk '{print $1}' | grep -v 'PROFILE' | xargs colima delete --force -p || true")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to delete " + k8Provider + " cause: ", err)
    }

    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "colima prune --force || true")
}

func listExistingMinikube() {
    cmd := exec.Command( "/bin/sh", "-c", "minikube profile list | awk '{print $1}' | grep -v 'PROFILE' ")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to List " + k8Provider + " cause:  ", err)
    }
    color.Notice.Tips("List of existing " + k8Provider + " installs")
    color.Notice.Println(string(out))
}

func listExistingColima() {
    cmd := exec.Command( "/bin/sh", "-c", "colima list | awk '{print $1}' | grep -v 'PROFILE' ")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to List " + k8Provider + " cause:  ", err)
    }
    color.Notice.Tips("List of existing " + k8Provider + " installs")
    color.Notice.Println(string(out))
}

func applyNginxIngress() {
    cmd := exec.Command("/bin/sh", "-c", "helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx --force-update")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to add Nginx helm repo. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "helm repo update")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to update Nginx helm repo. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "helm install -n ingress-nginx --create-namespace --set controller.hostNetwork=true --set controller.watchIngressWithoutClass=true ingress-nginx ingress-nginx/ingress-nginx")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to install Nginx. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "kubectl get ingress nginx")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to check Nginx ingress, cause: : ", err)
    }
    color.Notice.Println(string(out))
}
func applyTraefik() {
    cmd := exec.Command("/bin/sh", "-c", "helm repo add traefik https://helm.traefik.io/traefik --force-update")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println(fmt.Sprint(err) + ": " + string(out))
        color.Error.Println("Failed to add Traefik helm repo. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "helm repo update")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to update Traefik helm repo. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "helm install traefik traefik/traefik --namespace=traefik --values=./parsed/traefik-values.yaml --set version=23.0.1")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println(fmt.Sprint(err) + ": " + string(out))
        color.Error.Println("Failed to install Traefik. Cause:  ", err)
    }
    color.Info.Println(string(out))
}

func applyTraefikCerts() {
    cmd := exec.Command("/bin/sh", "-c", "kubectl apply -f ./parsed/traefik-certs.yaml")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to add traefik certs. Cause: ", err)
    }
    color.Notice.Println(string(out))
}

func applyDashboard() {
    cmd := exec.Command("/bin/sh", "-c", "kubectl rollout status -w deployment/kubernetes-dashboard -n kubernetes-dashboard || true")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println(fmt.Sprint(err) + ": " + string(out))
        color.Error.Println("Failed to check k8s dashboard rollout. Cause: ", err)
    }
    color.Notice.Println(string(out))
    cmd = exec.Command("/bin/sh", "-c", "kubectl apply -f ./parsed/dashboard.yaml -n kubernetes-dashboard")
    out, err = cmd.Output()
    if err != nil {
        color.Error.Println("Failed to apply k8s dashboard ingress. Cause:  ", err)
    }
    color.Info.Println(string(out))
}

func installCertManager() {
    cmd := exec.Command( "/bin/sh", "-c", "helm repo add jetstack https://charts.jetstack.io;helm repo update;helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.6.3 --set installCRDs=true")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Failed to add Cert manager helm repo. Cause: ", err)
    }

    color.Info.Println(string(out))
}
func installCerts() {
    color.Info.Tips("Installing certs")
    command := []string{
        "./scripts/certs-install.sh",
    }

    execute("./scripts/certs-install.sh", command)
}
func prepDependencies() {
    var brewInstalls = [][]string{
        {"install", "yq"},
        {"install", "jq"},
        {"install", "stern"},
        {"install", "minikube"},
        {"install", "kubectl"},
        {"install", "docker"},
        {"install", "colima"},
        {"install", "cmctl"},
        {"install", "mkcert"},
        {"install", "openLens"},
        {"install", "kubectx"},
        {"install", "k9s"},
        {"install", "helm"},
        {"install","tilt-dev/tap/tilt"},
        {"install", "tilt-dev/tap/ctlptl"},
    }
   color.Info.Tips("Installing dependencies")
    for i := 0; i < len(brewInstalls); i++  {

        cmd := exec.Command("brew", brewInstalls[i]...)
        if errors.Is(cmd.Err, exec.ErrDot) {
        	cmd.Err = nil
        }
        if err := cmd.Run(); err != nil {
        	log.Fatal(err)
        } else {
            color.Notice.Println("Completed: brew",strings.Join(brewInstalls[i]," "))
        }
    }
    color.Info.Println("Installed dependencies")
}

func setupHosts(){
    cmd := exec.Command( "/bin/sh", "-c", "sudo -- sh -c -e \"echo '127.0.0.1   test.localhost' >> /etc/hosts\"")
    out, err := cmd.Output()
    if err != nil {
        color.Error.Println("Could not run command: ", err)
    }
    color.Notice.Println(string(out))
}

// func prepTemplates() error {
//     const templateFile = "templates/template.yaml"
//     const dataFile = "templates/values.yaml"
//     const outputFile = "parsed.yaml"
//     const targetDir = "parsed"
//     if err := parser.Parse(templateFile, dataFile, outputFile,targetDir); err != nil {
//         return err
//     }
//     fmt.Printf("File %s was generated.\n", outputFile)
//     return nil
// }

func buildr(values map[string]string) error {
    const templateFile = "templates/resolver.yaml"
    const outputFile = "parsedResolver.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}
func buildClusterIssuer(values map[string]string) error {
    const templateFile = "templates/cluster-issuer.yaml"
    const outputFile = "cluster-issuer.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}

func buildTestApp(values map[string]string) error {
    const templateFile = "templates/test-app.yaml"
    const outputFile = "test-app.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}

func buildTraefikValues(values map[string]string) error {
    const templateFile = "templates/traefik-values.yaml"
    const outputFile = "traefik-values.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}

func buildTraefikCerts(values map[string]string) error {
    const templateFile = "templates/traefik-certs.yaml"
    const outputFile = "traefik-certs.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}
func buildResolver(values map[string]string) error {
    const templateFile = "templates/resolver.yaml"
    const outputFile = "parsedResolver.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}

func buildDashboardIngress(values map[string]string) error {
    const templateFile = "templates/dashboard.yaml"
    const outputFile = "dashboard.yaml"
    const targetDir = "parsed"
    if err := parser.ParseValues(templateFile, values, outputFile,targetDir); err != nil {
        return err
    }
    fmt.Printf("File %s was generated.\n", outputFile)
    return nil
}

