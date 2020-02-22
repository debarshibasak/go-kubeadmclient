package kubeadmclient

import (
	"encoding/json"
	"fmt"
	"github.com/debarshibasak/kubekray/sshclient"
	"time"
)

type MasterNode struct {
	*Node
}

func NewMasterNode(username string,
	ipOrHost string,
	privateKeyLocation string) *MasterNode {

	return &MasterNode{&Node{
		username:           username,
		ipOrHost:           ipOrHost,
		privateKeyLocation: privateKeyLocation,
	},
	}
}

func NewMasterNodes(username string,
	ipOrHost []string,
	osType string,
	privateKeyLocation string) []MasterNode {

	var masterNodes []MasterNode

	for _, ip := range ipOrHost {
		masterNodes = append(masterNodes, MasterNode{&Node{
			username:           username,
			ipOrHost:           ip,
			osType:             osType,
			privateKeyLocation: privateKeyLocation,
		}})
	}

	return masterNodes
}

func (n *MasterNode) GetToken() (string, error) {

	sh := sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
	}

	out, err := sh.Collect("sudo kubeadm token list -o json")
	if err != nil {
		return "", err
	}

	c := make(map[string]interface{})

	err = json.Unmarshal([]byte(out), &c)
	if err != nil {
		return "",  err
	}

	return c["token"].(string), nil
}

func (n *MasterNode) GetKubeConfig() (string, error) {
	return n.sshClient().Collect("sudo cat /etc/kubernetes/admin.conf")
}

func (n *MasterNode) GetJoinCommand() (string, error) {
	return n.sshClient().Collect("sudo kubeadm token create --print-join-command")
}

func (n *MasterNode) Install() error {

	osType := n.determineOS()

	fmt.Println("os determined "+osType)

	var cmds []string

	if osType == "ubuntu" {
		cmds = []string{
			"sudo apt-get update",
			"sudo apt-get install -y iptables arptables ebtables",
			//"sudo update-alternatives --set iptables /usr/sbin/iptables-legacy",
			//"sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy",
			//"sudo update-alternatives --set arptables /usr/sbin/arptables-legacy",
			//"sudo update-alternatives --set ebtables /usr/sbin/ebtables-legacy",
			"sudo apt-get update && sudo apt-get install -y apt-transport-https curl",
			"curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -",
			`cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
`,
			"sudo apt-get update",
			"sudo apt-get install -y kubelet kubeadm kubectl",
			"sudo apt-mark hold kubelet kubeadm kubectl",
			"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
			`sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"`,
			"sudo apt update",
			"apt-cache policy docker-ce",
			"sudo apt install docker-ce -y",
			"sudo usermod -aG docker ${USER}",
			"sudo swapoff -a",
			`sudo sed -i "/ swap / s/^\(.*\)$/#\1/g" /etc/fstab`,
			"sudo sysctl net.bridge.bridge-nf-call-iptables=1",
		}
	}

	err := n.sshClient().Run(cmds)
	if err != nil {
		return err
	}

	return n.sshClientWithTimeout(20*time.Minute).Run([]string{
			"sudo kubeadm init --pod-network-cidr=192.178.0.0/16 --service-cidr=192.178.168.0/16",
	})
}