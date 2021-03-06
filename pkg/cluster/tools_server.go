package cluster

import (
	"fmt"
)

func ToolsServerCreate(hcloudToken string) error {
	fmt.Println("executing \"ToolsServerCreate()\"")
	// backgroundCtx := context.Background()
	// client := hcloud.NewClient(hcloud.WithToken(hcloudToken))

	//create(initializedProviders[getProviderSting(clusterCfg.toolsServer)])

	return nil
}

func CreateToolsServer(clusterCfg Config) {
	// ### setup toolsServer
	// - check if toolsServer is wanted (means is described in cluster.yaml)
	//   - exit if yes (bec. it's optinal)
	// - check provider connection + authentication
	// - gather_facts
	// - check if toolsServer already exists
	//   - exit if yes
	// - provision toolsServer
	// - install (cli) tools
	//   - 100%
	//     - nano
	//     - vi(m)
	//     - ssh
	//     - openssl
	//     - kubectl
	//   - optional
	//     - git
	//     - helm

}
