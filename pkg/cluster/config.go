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
	Static []struct{} `yaml:"static"`
	Groups struct {
		Master masterServers `yaml:"master"`
		Worker workerServers `yaml:"worker"`
	} `yaml:"groups"`
}

type masterServers []server

type workerServers []server

type server struct {
	BaseName                        string                  `yaml:"baseName"`
	Scale                           int                     `yaml:"scale"`
	EnablePossibilityToProxyKubectl bool                    `yaml:"enablePossibilityToProxyKubectl"`
	ProviderMachineTemplate         providerMachineTemplate `yaml:"providerMachineTemplate"`
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
