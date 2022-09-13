package cluster

import (
	"fmt"

	"github.com/pandorasnox/kubelife/pkg/ssh"
)

// Config for the cluster
type Config struct {
	Version string  `yaml:"version"`
	Cluster Cluster `yaml:"cluster"`
}

// Cluster related information
type Cluster struct {
	Name              string            `yaml:"name"`
	SSHAuthorizedKeys sshAuthorizedKeys `yaml:"ssh_authorized_keys"`
	Nodes             servers           `yaml:"nodes"`
}

type servers struct {
	Static struct {
		ToolsServer toolsServer    `yaml:"toolsServer"`
		Worker      []staticServer `yaml:"worker"`
	} `yaml:"static"`
	Scalable struct {
		ControlPlane scalableControlPlaneServers `yaml:"controlPlane"`
		Worker       scalableWorkerServers       `yaml:"worker"`
	} `yaml:"scalable"`
}

type sshAuthorizedKeys []ssh.PubKeyData

func (sshKeys sshAuthorizedKeys) NameList() []string {
	list := []string{}
	for _, v := range sshKeys {
		list = append(list, v.Name)
	}

	return list
}

type scalableControlPlaneServers []scalableServer

type scalableWorkerServers []scalableServer

type staticServer struct {
	NameAddition            string                  `yaml:"nameAddition"`
	ProviderMachineTemplate providerMachineTemplate `yaml:"providerMachineTemplate"`
}

type scalableServer struct {
	NameAddition            string                  `yaml:"nameAddition"`
	Scale                   int                     `yaml:"scale"`
	ProviderMachineTemplate providerMachineTemplate `yaml:"providerMachineTemplate"`
}

type toolsServer struct {
	// Name                    string                  `yaml:"name"`
	ProviderMachineTemplate providerMachineTemplate `yaml:"providerMachineTemplate"`
}

type providerMachineTemplate struct {
	HetznerCloud hetznerCloudMachine `yaml:"hetznerCloud"` //pointer? / yaml anotaion when ignore when not set "omitempty"
	Digitalocean digitaloceanMachine `yaml:"digitalocean"`
}

type hetznerCloudMachine struct {
	ServerType string `yaml:"serverType"`
	Image      struct {
		Name string `yaml:"name"`
	} `yaml:"image"`
	Location         string `yaml:"location"`
	AdditionalLabels Labels `yaml:"additionalLabels"`
}

type digitaloceanMachine struct {
	ServerType string `yaml:"serverType"`
}

type Labels map[string]string

func MergeLabels(left Labels, right Labels) (out Labels, warnings []string) {
	out = Labels{}

	for k, v := range left {
		out[k] = v
	}

	for k, v := range right {
		if _, ok := out[k]; ok {
			warnMsg := fmt.Sprintf("SKIPPING THIS ACTION: you try to add a label that already exists: \"%s: %s\" ({key: \"%s\", value: \"%s\"})", k, v, k, v)
			warnings = append(warnings, warnMsg)
			continue
		}

		out[k] = v
	}

	return out, warnings
}
