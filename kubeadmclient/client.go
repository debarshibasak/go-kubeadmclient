package kubeadmclient

import "github.com/debarshibasak/kubekray/sshclient"

type Node struct {
	username           string
	ipOrHost           string
	osType             string
	privateKeyLocation string
}

type MasterNode struct {
	*Node
}

type WorkerNode struct {
	*Node
	token string
}

func NewWorkerNode(username string,
	ipOrHost string,
	osType string,
	privateKeyLocation string) *WorkerNode {

	return &WorkerNode{
		Node: &Node{
			username:           username,
			ipOrHost:           ipOrHost,
			osType:             osType,
			privateKeyLocation: privateKeyLocation,
		},
	}
}

func NewMasterNode(username string,
	ipOrHost string,
	osType string,
	privateKeyLocation string) *MasterNode {

	return &MasterNode{&Node{
		username:           username,
		ipOrHost:           ipOrHost,
		osType:             osType,
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

	cmds := []string{
		"sudo kubeadm token list -o json",
	}

	if err := sh.Run(cmds); err != nil {
		return "", err
	}

	return "", nil
}

func (n *MasterNode) GetKubeConfig() (string, error) {
	return "", nil
}

func (n *MasterNode) Install() error {

	sh := sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
	}

	cmds := []string{
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
		"sudo kubeadm init",
	}

	return sh.Run(cmds)
}

func (n *WorkerNode) SetToken(token string) {
	n.token = token
}

func (n *WorkerNode) Install() error {

	sh := sshclient.SshConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
	}

	cmds := []string{
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
		n.token,
	}

	return sh.Run(cmds)
}
