package kubeadmclient

import (
	"github.com/debarshibasak/go-kubeadminclient/sshclient"
	"strings"
	"time"
)

type Node struct {
	username           string
	ipOrHost           string
	osType             string
	privateKeyLocation string
	verboseMode        bool
	clientID           string
}

func (n *Node) determineOS() string {

	client := n.sshClient()
	out, err := client.Collect("uname -a")
	if err != nil {
		return "n/a"
	}

	if strings.Contains(out, "Ubuntu") {
		return "ubuntu"
	}

	if err := client.Run([]string{"ls /etc/centos-release"}); err == nil {
		return "centos"
	}

	if err := client.Run([]string{"ls /etc/redhat-release"}); err == nil {
		return "redhat"
	}

	return "n/a"
}

func (n *Node) sshClient() *sshclient.SshConnection {
	return &sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
		VerboseMode: n.verboseMode,
		ClientID:    n.clientID,
	}
}

func (n *Node) sshClientWithTimeout(duration time.Duration) *sshclient.SshConnection {
	return &sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
		VerboseMode: n.verboseMode,
		Timeout:     duration,
	}
}
