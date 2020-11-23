/*
Copyright 2020 The Kubelife Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pandorasnox/kubelife/pkg/cluster"
	"gopkg.in/yaml.v2"

	"github.com/pandorasnox/kubelife/pkg/hetzner"
	"github.com/pandorasnox/kubelife/pkg/ssh"
	cli "github.com/urfave/cli/v2"
)

func main() {
	// user := flag.String("user", "", "remote server login user")
	// addr := flag.String("addr", "", "remote server address (ip/dns)")
	// flag.Parse()
	// _ = user
	// _ = addr

	// osSetup(*user, *addr)

	var err error

	clusterCfg, err := loadClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = clusterCfg

	app := &cli.App{
		Name:  "Kubelife",
		Usage: "usage: todo",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "hetzner",
				// Aliases: []string{"h"},
				Usage: "all hetzner cloud related commands",
				Subcommands: []*cli.Command{
					{
						Name: "status",
						// Aliases: []string{"s"},
						Usage: "prints the hetzner status to std.out",
						Action: func(c *cli.Context) error {
							err := hetzner.Status(os.Getenv("HCLOUD_TOKEN"))
							if err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name: "server",
						// Aliases: []string{"s"},
						Usage: "commands related to server",
						Subcommands: []*cli.Command{
							{
								Name: "status",
								// Aliases: []string{"s"},
								Usage: "prints the hetzner status to std.out",
								Action: func(c *cli.Context) error {
									err := hetzner.Status(os.Getenv("HCLOUD_TOKEN"))
									if err != nil {
										return err
									}

									return nil
								},
							},
							{
								Name: "create",
								// Aliases: []string{"s"},
								Usage: "creates a new hetzner cloud vm",
								Action: func(c *cli.Context) error {
									err := hetzner.Create(os.Getenv("HCLOUD_TOKEN"))
									if err != nil {
										return err
									}

									return nil
								},
							},
							{
								Name: "delete",
								// Aliases: []string{"s"},
								Subcommands: []*cli.Command{
									{
										Name: "all",
										// Aliases: []string{"s"},
										Usage: "deletes all vms",
										Action: func(c *cli.Context) error {
											err := hetzner.DeleteAll(os.Getenv("HCLOUD_TOKEN"))
											if err != nil {
												return err
											}

											return nil
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadClusterConfig() (cluster.Config, error) {
	// load file
	f, err := os.Open("cluster.yml")
	if err != nil {
		log.Fatalf("open config: %v", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("stat config: %v", err)
	}
	if fi.Size() == 0 {
		log.Fatalf("cluster config is empt")
	}

	// parse/decode file
	decodedCfg, err := decodeClusterConfig(f)
	if err != nil {
		log.Fatalf("parse / decode config: %v", err)
	}

	//maybe add defaults

	// validate / sanity check file
	//todo

	return decodedCfg, nil
}

// LoadConfig load the config from the reader.
func decodeClusterConfig(r io.Reader) (cluster.Config, error) {
	d := yaml.NewDecoder(r)
	d.SetStrict(true)

	cfg := cluster.Config{}

	err := d.Decode(&cfg)
	if err != nil {
		return cluster.Config{}, err
	}

	return cfg, nil
}

func osSetup(user string, remoteAddrs string) {
	ssh, err := ssh.New(user, remoteAddrs, ssh.AgentAuth())
	if err != nil {
		log.Fatalf("foo %s", err)
	}
	defer ssh.Close()

	log.Println("update system")
	_, err = ssh.Exec("apt-get update && apt-get upgrade -y")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

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
			log.Fatalf("bar %s", err)
		}
	}

	log.Println("install & enable docker")
	_, err = ssh.Exec("apt-get install -y docker.io && systemctl enable docker.service")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

	log.Println("add kubernetes packae list & update")
	k8sPkgCmd := "apt-get update && apt-get install -y apt-transport-https curl"
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -")
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "echo deb https://apt.kubernetes.io/ kubernetes-xenial main > /etc/apt/sources.list.d/kubernetes.list")
	k8sPkgCmd = fmt.Sprintf("%s && %s", k8sPkgCmd, "apt-get update")
	_, err = ssh.Exec(k8sPkgCmd)
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

	log.Println("add k8s tooling")
	_, err = ssh.Exec("apt-get install -y kubelet=1.18.6-00 kubeadm=1.18.6-00 kubectl=1.18.6-00")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

	log.Println("hold k8s tooling")
	_, err = ssh.Exec("apt-mark hold kubelet kubeadm kubectl")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

	log.Println("enable iptables")
	_, err = ssh.Exec("modprobe br_netfilter && sysctl net.bridge.bridge-nf-call-iptables=1 && sysctl net.bridge.bridge-nf-call-ip6tables=1")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

	log.Println("disable swap")
	_, err = ssh.Exec("swapoff -a && sed -i '/ swap / s/^/#/' /etc/fstab")
	if err != nil {
		log.Fatalf("bar %s", err)
	}
	// fmt.Printf("%s\n", out)

}
