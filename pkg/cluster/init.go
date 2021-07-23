package cluster

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/pandorasnox/kubelife/pkg/hetzner"
	"github.com/pandorasnox/kubelife/pkg/ssh"
	log "github.com/sirupsen/logrus"
)

func Init(ccfg Config, hcloud_token string) error {
	var err error

	err = addSSHKeysToProvider(hcloud_token, "hetzner_cloud", ccfg.Cluster.Nodes.SSHAuthorizedKeys)
	if err != nil {
		return err
	}

	err = initToolsServer(ccfg, hcloud_token)
	if err != nil {
		return fmt.Errorf("couldn't initiate toolsServer: %s", err)
	}

	return nil
}

func addSSHKeysToProvider(hcloud_token string, provider string, sshKeys []ssh.PubKeyData) error {
	if provider == "hetzner_cloud" {
		log.Info("add ssh public keys to provider hetzner_cloud")
		err := hetzner.CreateSSHKeys(hcloud_token, sshKeys)
		if err != nil {
			return err
		}
	}

	return nil
}

func initToolsServer(ccfg Config, hcloud_token string) error {
	var err error

	cToolsServer := ccfg.Cluster.Nodes.ToolsServer
	if reflect.ValueOf(cToolsServer).IsZero() {
		log.Println("skip toolsServer initialisation, given empty value(s)")
		return nil
	}

	if reflect.ValueOf(cToolsServer.ProviderMachineTemplate).IsZero() {
		msg := "for the toolsServer you need to provide a concrete ProviderMachineTemplate"
		return errors.New(msg)
	}

	v := reflect.ValueOf(cToolsServer.ProviderMachineTemplate)
	countNonEmpty := 0
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			countNonEmpty++
		}
	}

	if countNonEmpty > 1 {
		return errors.New("for the toolsServer you provided more than 1 ProviderMachineTemplate")
	}

	provider, err := extractFirstFound(v)
	if err != nil {
		return fmt.Errorf("couldn't extract ProviderMachineTemplate: %s", err)
	}

	toolsServerName := fmt.Sprintf("%s%s", ccfg.Cluster.Name, "-clustertools")

	switch provider.Interface().(type) {
	case hetznerCloudMachine:
		hcm, _ := provider.Interface().(hetznerCloudMachine)

		creatorLabels := Labels{
			"kubelife_owner":   "kubelife",
			"kubelife_creator": "kubelife",
		}

		mergedLabels, warnings := MergeLabels(creatorLabels, hcm.AdditionalLabels)
		if len(warnings) != 0 {
			for _, v := range warnings {
				log.Warn(v)
			}
		}

		allHSSHKeys, err := hetzner.SSHKeys(hcloud_token)
		if err != nil {
			return err
		}

		sshKeysForMachine := hetzner.FilterSSHKeysByNameList(
			allHSSHKeys, ccfg.Cluster.Nodes.SSHAuthorizedKeys.NameList(),
		)

		//map hetznerCloudMachine => to => hcloud.ServerCreateOpts
		hcScOps := hcloud.ServerCreateOpts{
			Name: toolsServerName,
			ServerType: &hcloud.ServerType{
				Name: hcm.ServerType,
			},
			Image: &hcloud.Image{
				Name: hcm.Image.Name,
			},
			Labels:  mergedLabels,
			SSHKeys: sshKeysForMachine,
		}

		//if toolsServer already exists, skip (add flag to force re-creation)

		err = hetzner.Create(hcloud_token, hcScOps, toolsServerName)
		if err != nil {
			return fmt.Errorf("couldn't create toolsServer as a hetznerCloudMachine: %s", err)
		}

		hToolsServer, err := hetzner.WaitForServerRunning(hcloud_token, toolsServerName, 35*time.Second)
		if err != nil {
			return fmt.Errorf("waiting for toolsServer is running failed: %s", err)
		}

		// wait for ssh access works
		err = waitForSSH("root", hToolsServer.PublicNet.IPv4.IP.String(), 15*time.Second)
		if err != nil {
			return fmt.Errorf("waiting for toolsServer ssh access failed: %s", err)
		}

		// os install tools / packages
		err = installPackagesForToolsServer("root", hToolsServer.PublicNet.IPv4.IP.String())
		if err != nil {
			return fmt.Errorf("couldn't install os packages for toolsServer: %s", err)
		}
	}

	return nil
}

func extractFirstFound(v reflect.Value) (reflect.Value, error) {
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			return v.Field(i), nil
		}
	}

	return reflect.Value{}, errors.New("coudn't extract/found even one")
}

func waitForSSH(user string, remoteAddrs string, timeout time.Duration) error {
	log.Infof("waiting for ssh access for \"%s\", timeout is set to \"%s\"", remoteAddrs, timeout.String())

	if timeout <= 0 {
		return errors.New("timeout needs to be larger than 0")
	}

	start := time.Now()
	end := start.Add(timeout)

	var lastErr error
	for {
		now := time.Now()
		if now.After(end) {
			return fmt.Errorf("reached timeout of '%s' seconds, last err: %s", timeout, lastErr)
		}

		ssh, err := ssh.New(user, remoteAddrs, ssh.AgentAuth())
		if err != nil {
			lastErr = fmt.Errorf("couldn't create ssh client: %s", err)
		}

		if err == nil {
			ssh.Close()
			break
		}

		time.Sleep(1 * time.Second)
	}

	log.Infof("ssh access is now available for \"%s\"", remoteAddrs)

	return nil
}

func installPackagesForToolsServer(user string, remoteAddrs string) error {
	ssh, err := ssh.New(user, remoteAddrs, ssh.AgentAuth())
	if err != nil {
		return fmt.Errorf("couldn't create ssh client: %s", err)
	}
	defer ssh.Close()

	log.Println("update system")
	_, err = ssh.Exec("apt-get update && apt-get upgrade -y")
	if err != nil {
		return fmt.Errorf("couldn't update system: %s", err)
	}

	log.Println("install os packages")
	packages := []string{
		"htop",
		"iotop",
		"atop",
		"nload",
		"sysstat",
		"smartmontools",
		// "docker.io",
		"ethtool",
		"socat",
		"dnsutils",
		"bash-completion",
		"bsdmainutils",
	}
	for _, pkg := range packages {
		_, err = ssh.Exec(fmt.Sprintf("apt install -y %s", pkg))
		if err != nil {
			return fmt.Errorf("couldn't install os package \"%s\": %s", pkg, err)
		}
	}

	log.Println("install & enable docker")
	_, err = ssh.Exec("apt-get install -y docker.io && systemctl enable docker.service")
	if err != nil {
		return fmt.Errorf("couldn't install docker: %s", err)
	}

	log.Println("enable docker")
	_, err = ssh.Exec("systemctl enable docker.service")
	if err != nil {
		return fmt.Errorf("couldn't enable docker: %s", err)
	}

	log.Println("add kubernetes packae list & update")
	k8sPkgCmd := "apt-get update && apt-get install -y apt-transport-https curl"
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -")
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "echo deb https://apt.kubernetes.io/ kubernetes-xenial main > /etc/apt/sources.list.d/kubernetes.list")
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "apt-get update")
	_, err = ssh.Exec(k8sPkgCmd)
	if err != nil {
		return fmt.Errorf("couldn't add kubernetes packae list & update: %s", err)
	}

	log.Println("install kubectl")
	_, err = ssh.Exec("apt-get install -y kubectl=1.18.6-00")
	if err != nil {
		log.Fatalf("couldn't install kubectl: %s", err)
	}

	log.Println("hold k8s tooling")
	_, err = ssh.Exec("apt-mark hold kubectl")
	if err != nil {
		log.Fatalf("bar %s", err)
	}

	log.Println("disable swap")
	_, err = ssh.Exec("swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab")
	if err != nil {
		log.Fatalf("bar %s", err)
	}

	return nil
}

//plan
////gather info
////create plan
//applay
////exec plan

//reconsile (currentState, desiredState) error
//// <- method or gets hetzner client

// interface cloud provide
// cloudProvider.reconcile(desiredState)

// how to notice vms ?
