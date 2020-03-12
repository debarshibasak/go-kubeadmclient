package kubeadmclient

import (
	"fmt"
	"log"
	"strings"
	"time"

	osType "github.com/debarshibasak/go-kubeadmclient/kubeadmclient/ostype"

	"github.com/debarshibasak/go-kubeadmclient/sshclient"
)

type Node struct {
	username           string
	ipOrHost           string
	osType             string
	privateKeyLocation string
	verboseMode        bool
	clientID           string
}

func (n *Node) String() string {
	return fmt.Sprintf("ip=%v username=%v key=%v", n.ipOrHost, n.username, n.privateKeyLocation)
}

func (n *Node) determineOS() osType.OsType {

	client := n.sshClient()
	out, err := client.Collect("uname -a")
	if err != nil {
		log.Println(string(out))
		return nil
	}

	if strings.Contains(out, "Ubuntu") {
		return &osType.Ubuntu{}
	}

	if err := client.Run([]string{"ls /etc/centos-release"}); err == nil {
		return &osType.Centos{}
	}

	if err := client.Run([]string{"ls /etc/redhat-release"}); err == nil {
		return &osType.Centos{}
	}

	return nil
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
