package hetzner

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/pandorasnox/kubelife/pkg/ssh"
)

// Status prints an overview status of hetzner cloud to std.out
func Status(token string) error {
	backgroundCtx := context.Background()
	client := hcloud.NewClient(hcloud.WithToken(token))

	actions, err := Actions(client, backgroundCtx)
	if err != nil {
		return fmt.Errorf("error retrieving actions: %s", err)
	}
	log.Printf("actions: %v\n", actions)

	serverTypes, err := client.ServerType.All(backgroundCtx)
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return fmt.Errorf("error retrieving server: %s", err)
	}
	// _ = serverTypes

	for _, serverTypePointer := range serverTypes {
		log.Printf("server type pointer . name: %v\n", serverTypePointer.Name)
		fmt.Println("")
	}
	log.Printf("server types: %v\n", serverTypes)

	hSSHKeyNames := []string{}
	hSSHKeys, err := client.SSHKey.All(backgroundCtx)
	for _, hSSHKey := range hSSHKeys {
		hSSHKeyNames = append(hSSHKeyNames, hSSHKey.Name)
	}

	log.Printf("ssh keys names: %v\n", hSSHKeyNames)

	servers, err := Servers(client, backgroundCtx)
	if err != nil {
		return fmt.Errorf("error retrieving server: %s", err)
	}
	log.Printf("servers: %v\n", servers)

	for _, server := range servers {
		fmt.Printf("server name: \"%s\", status: \"%s\"", server.Name, server.Status)
	}

	return nil
}

func Actions(client *hcloud.Client, ctx context.Context) ([]*hcloud.Action, error) {
	actionOpts := hcloud.ActionListOpts{}
	actions, _, err := client.Action.List(ctx, actionOpts)
	if err != nil {
		return []*hcloud.Action{}, err
	}

	return actions, nil
}

func Servers(client *hcloud.Client, ctx context.Context) ([]*hcloud.Server, error) {
	servers, err := client.Server.All(context.Background())
	if err != nil {
		// log.Errorf("%s", err)
		// log.Fatalf("error retrieving server: %s\n", err)
		return []*hcloud.Server{}, fmt.Errorf("error retrieving server: %s", err)
	}

	return servers, nil
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

func Create(token string, hcScOps hcloud.ServerCreateOpts, serverName string) error {
	startAfterCreate := true
	hcScOps.StartAfterCreate = &startAfterCreate

	automount := false
	hcScOps.Automount = &automount

	location := ""
	if location != "" {
		hcScOps.Location = &hcloud.Location{Name: location}
	}

	client := hcloud.NewClient(hcloud.WithToken(token))

	_, _, err := client.Server.Create(context.Background(), hcScOps)
	if err != nil {
		return fmt.Errorf("couldn't create  server: %s", err)
	}

	log.Printf("created server: \"%s\"\n", serverName)

	return nil
}

// Create a new opinionated hetzner cloud server
func CreateSingle(token string) error {
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

func WaitForServerRunning(token string, serverName string, timeoutSeconds time.Duration) error {
	log.Infof("waiting for server \"%s\" is running", serverName)

	backgroundCtx := context.Background()
	client := hcloud.NewClient(hcloud.WithToken(token))

	if timeoutSeconds <= 0 {
		return errors.New("seconds needs to be larger than 0")
	}

	start := time.Now()
	end := start.Add(timeoutSeconds * time.Second)

	for {
		now := time.Now()
		if !inTimeSpan(start, end, now) {
			return fmt.Errorf("reached timeout of '%s' seconds", timeoutSeconds)
		}

		server, _, err := client.Server.GetByName(backgroundCtx, serverName)
		if err != nil {
			return fmt.Errorf("error retrieving server with name '%s': %s", serverName, err)
		}

		if server.Status == hcloud.ServerStatusRunning {
			break
		}

		time.Sleep(1 * time.Second)
	}

	log.Infof("server \"%s\" is now running", serverName)

	return nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
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

func CreateSSHKeys(token string, sshKeys []ssh.PubKeyData) error {
	client := hcloud.NewClient(hcloud.WithToken(token))

	for _, sshKey := range sshKeys {
		sshOps := hcloud.SSHKeyCreateOpts{
			Name:      sshKey.Name,
			PublicKey: sshKey.PublicKey,
			// Labels    map[string]string
		}

		_, _, err := client.SSHKey.Create(context.Background(), sshOps)
		if err != nil {
			return fmt.Errorf("couldn't create ssh key with name \"%s\": %s", sshKey.Name, err)
		}
	}

	return nil
}
