package cluster

import (
	"errors"
	"fmt"
	"time"

	"github.com/pandorasnox/kubelife/pkg/environment"
	"github.com/pandorasnox/kubelife/pkg/hetzner"
	"github.com/pandorasnox/kubelife/pkg/ssh"
	log "github.com/sirupsen/logrus"
)

func Init(ccfg Config, env environment.Config) error {
	var err error

	// todo: run this only based on provider found in cluster.yaml, not hardcoded
	err = addSSHKeysToProvider(env.HcloudToken, "hetzner_cloud", ccfg.Cluster.SSHAuthorizedKeys)
	if err != nil {
		return err
	}

	err = provisionToolsServer(ccfg, env)
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

//plan
////gather info
////create plan
//applay
////exec plan

//reconsile (currentState, desiredState) error
//// <- method or gets hetzner client

// interface cloud provide
// cloudProvider.reconcile(desiredState)
// cloudProvider.Provision(ccfg, envCfg)

// how to notice vms ?

// cloudProvider.Provision(ccfg, envCfg)
// install os packages
