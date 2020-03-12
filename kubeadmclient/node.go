package kubeadmclient

import (
	"fmt"
	"log"
	"strings"
	"time"

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

type OsType string

const (
	Ubuntu    OsType = "UBUNTU"
	Centos    OsType = "CENTOS"
	RedHat    OsType = "REDHAT"
	UnknownOS OsType = "n/a"
)

func (n *Node) String() string {
	return fmt.Sprintf("ip=%v username=%v key=%v", n.ipOrHost, n.username, n.privateKeyLocation)
}

func (n *Node) determineOS() OsType {

	client := n.sshClient()
	out, err := client.Collect("uname -a")
	if err != nil {
		log.Println(string(out))
		return UnknownOS
	}

	if strings.Contains(out, "Ubuntu") {
		return Ubuntu
	}

	if err := client.Run([]string{"ls /etc/centos-release"}); err == nil {
		return Centos
	}

	if err := client.Run([]string{"ls /etc/redhat-release"}); err == nil {
		return RedHat
	}

	return UnknownOS
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
