package kubeadmclient

import (
	"github.com/debarshibasak/kubekray/sshclient"
	"time"
)

type WorkerNode struct {
	*Node
}

func NewWorkerNode(username string,
	ipOrHost string,
	privateKeyLocation string) *WorkerNode {

	return &WorkerNode{
		Node: &Node{
			username:           username,
			ipOrHost:           ipOrHost,
			privateKeyLocation: privateKeyLocation,
		},
	}
}

func (n *WorkerNode) Install(joinCommand string) error {

	sh := sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
	}

	cmds := []string{
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

	}

	if err := sh.Run(cmds); err != nil {
		return err
	}

	return n.sshClientWithTimeout(30*time.Minute).Run([]string{
		"sudo "+joinCommand,
	})
}