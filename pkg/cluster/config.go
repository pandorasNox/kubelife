package cluster

// Config for the cluster
type Config struct {
	version string
	cluster Cluster
}

// Cluster related information
type Cluster struct {
	name  string
	nodes servers
}

// type servers []server
// type server struct{}

type servers struct {
	masters masterServers
	workers workerServers
}

type masterServers []masterServer
type masterServer struct{}

type workerServers []workerServer
type workerServer struct{}
