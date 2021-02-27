package cluster

// Config for the cluster
type Config struct {
	Version     string   `yaml:"version"`
	Cluster     Cluster  `yaml:"cluster"`
	ToolsServer struct{} `yaml:"toolsServer"`
}

// Cluster related information
type Cluster struct {
	Name  string  `yaml:"name"`
	Nodes servers `yaml:"nodes"`
}

type servers struct {
	SSHAuthorizedKeys []struct {
		Name      string `yaml:"name"`
		PublicKey string `yaml:"public_key"`
	} `yaml:"ssh_authorized_keys"`
	Static struct {
		Master []staticServer `yaml:"master"`
		Worker []staticServer `yaml:"worker"`
	} `yaml:"static"`
	ScalableGroups struct {
		Master scalableMasterServers `yaml:"master"`
		Worker scalableWorkerServers `yaml:"worker"`
	} `yaml:"scalableGroups"`
}

type scalableMasterServers []scalableServer

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

type providerMachineTemplate struct {
	HetznerCloud hetznerCloudMachineProvider `yaml:"hetznerCloud"`
}

type hetznerCloudMachineProvider struct {
	ServerType string `yaml:"serverType"`
	Image      struct {
		Name string `yaml:"name"`
	} `yaml:"image"`
	Location string            `yaml:"location"`
	Labels   map[string]string `yaml:"labels"`
}
