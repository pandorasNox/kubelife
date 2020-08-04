package ssh

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSH struct {
	Client *ssh.Client
}

// Connect to ssh and get client, the host public key must be in known hosts.
func New(user string, addr string, auth AuthFn) (ssh *SSH, err error) {
	return ssh, nil
}

type AuthFn func() (ssh.AuthMethod, error)

func AgentAuth() AuthFn {
	return func() (ssh.AuthMethod, error) {
		sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return nil, fmt.Errorf("could not dial/connect with ssh agent: %w", err)
		}
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), nil
	}
}
