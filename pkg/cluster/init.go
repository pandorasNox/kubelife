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

		//map hetznerCloudMachine => to => hcloud.ServerCreateOpts
		hcScOps := hcloud.ServerCreateOpts{
			Name: toolsServerName,
			ServerType: &hcloud.ServerType{
				Name: hcm.ServerType,
			},
			Image: &hcloud.Image{
				Name: hcm.Image.Name,
			},
			Labels: mergedLabels,
		}

		//if toolsServer already exists, skip (add flag to force re-creation)

		err := hetzner.Create(hcloud_token, hcScOps, toolsServerName)
		if err != nil {
			return fmt.Errorf("couldn't create toolsServer as a hetznerCloudMachine: %s", err)
		}

		err = hetzner.WaitForServerRunning(hcloud_token, toolsServerName, 5*time.Second)
		if err != nil {
			return fmt.Errorf("waiting for toolsServer is running failed: %s", err)
		}

		// wait for ssh access works

		// os install tools / packages
		err = installPackagesForToolsServer("user.name", "remote.address")
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

func installPackagesForToolsServer(user string, remoteAddrs string) error {
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
