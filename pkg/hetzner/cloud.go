package hetzner

import (
	"context"
	"fmt"
	"log"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// Status prints an overview status of hetzner cloud to std.out
func Status(token string) error {
	client := hcloud.NewClient(hcloud.WithToken(token))

	serverTypes, err := client.ServerType.All(context.Background())
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return fmt.Errorf("error retrieving server: %s", err)
	}
	_ = serverTypes

	// for serverType := range serverTypes {
	// 	log.Printf("server type: %v\n", serverType)
	// }
	// log.Printf("server types: %v\n", serverTypes)

	servers, err := client.Server.All(context.Background())
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return fmt.Errorf("error retrieving server: %s", err)
	}

	log.Printf("servers: %v\n", servers)

	return nil
}

func Locations(token string) error {
	client := hcloud.NewClient(hcloud.WithToken(token))

	locations, err := client.Location.All(context.Background())
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return fmt.Errorf("error retrieving locations: %s", err)
	}

	for _, l := range locations {
		log.Printf("location: %v\n", l)
	}

	return nil
}

// Create a new opinionated hetzner cloud server
func Create(token string) error {
	startAfterCreate := true
	automount := false
	opts := hcloud.ServerCreateOpts{
		Name: "test1",
		ServerType: &hcloud.ServerType{
			Name: "cx11",
		},
		Image: &hcloud.Image{
			Name: "ubuntu-20.04",
		},
		// SSHKeys          []*SSHKey
		// Location         *Location
		// Datacenter       *Datacenter // is discouraged
		// UserData         string
		StartAfterCreate: &startAfterCreate,
		Labels: map[string]string{
			"a": "b",
		},
		Automount: &automount,
		// Volumes          []*Volume
		// Networks         []*Network
	}
	location := ""
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}

	client := hcloud.NewClient(hcloud.WithToken(token))

	// stc := client.ServerType.All()

	// svrOps := &hcloud.ServerCreateOpts{
	// 	Name:       "test",
	// 	ServerType: &hcloud.ServerType{},
	// }
	scResult, res, err := client.Server.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("couldn't create  server: %s", err)
	}

	log.Printf("created server: %v\n", scResult)
	log.Printf("created server response: %v\n", res)

	return nil
}

func DeleteAll(token string) error {
	client := hcloud.NewClient(hcloud.WithToken(token))

	servers, err := client.Server.All(context.Background())
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return fmt.Errorf("error retrieving server: %s", err)
	}

	for i := range servers {
		server := servers[i]
		client.Server.Delete(context.Background(), server)
	}

	return nil
}
