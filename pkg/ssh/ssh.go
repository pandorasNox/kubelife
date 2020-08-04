package ssh

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

type SSH struct {
	Client *ssh.Client
}

// type Endpoint struct {
//     Host string
//     Port int
//     User string
// }
// func NewEndpoint(s string) *Endpoint {
//     endpoint := &Endpoint{
//         Host: s,
//     }
//     if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
//         endpoint.User = parts[0]
//         endpoint.Host = parts[1]
//     }
//     if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
//         endpoint.Host = parts[0]
//         endpoint.Port, _ = strconv.Atoi(parts[1])
//     }
//     return endpoint
// }

// Connect to ssh and get client, the host public key must be in known hosts.
func New(serverUser string, addr string, authFn AuthFn) (*SSH, error) {
	authMeth, err := authFn()
	if err != nil {
		return &SSH{}, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return &SSH{}, fmt.Errorf("could not create hostkeycallback function: %s", err)
	}

	knownHostPath := fmt.Sprintf("%s/%s", currentUser.HomeDir, ".ssh/known_hosts")

	hostKeyCallback, err := kh.New(knownHostPath)
	if err != nil {
		return &SSH{}, fmt.Errorf("could not create hostkeycallback function: %s", err)
	}

	config := &ssh.ClientConfig{
		User: serverUser,
		Auth: []ssh.AuthMethod{
			authMeth,
		},
		Timeout: 20 * time.Second,
		//see https://skarlso.github.io/2019/02/17/go-ssh-with-host-key-verification/
		HostKeyCallback: hostKeyCallback,
	}

	port := 22
	proto := "tcp"

	client, err := ssh.Dial(proto, fmt.Sprintf("%s:%d", addr, port), config)
	if err != nil {
		return &SSH{}, err
	}

	return &SSH{Client: client}, nil
}

func (s *SSH) Close() error {
	err := s.Client.Close()
	if err != nil {
		return err
	}

	return nil
}

// Exec runs a command on the remote server.
func (s *SSH) Exec(cmd string) (string, error) {
	sess, err := s.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create ssh session: %s", err)
	}
	defer sess.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	sess.Stdout = &stdoutBuf
	sess.Stderr = &stderrBuf

	err = sess.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("failed executing cmd (on remote): %s | stderr: %s", err, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}

type AuthFn func() (ssh.AuthMethod, error)

func AgentAuth() AuthFn {
	return func() (ssh.AuthMethod, error) {
		sshAgentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return nil, fmt.Errorf("Failed to open/connect SSH_AUTH_SOCK: %v", err)
		}
		return ssh.PublicKeysCallback(agent.NewClient(sshAgentConn).Signers), nil
	}
}
