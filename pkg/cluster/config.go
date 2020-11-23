package cluster

// Config for the cluster
type Config struct {
	Version string  `yaml:"version"`
	Cluster Cluster `yaml:"cluster"`
}

// Cluster related information
type Cluster struct {
	Name  string  `yaml:"name"`
	Nodes servers `yaml:"nodes"`
}

type servers struct {
	Master masterServers `yaml:"master"`
	Worker workerServers `yaml:"worker"`
}

type masterServers []server

type workerServers []server

type server struct {
	Name string `yaml:"name"`
}
