package kubeadmclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/centos"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/common"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/ubuntu"
	"github.com/debarshibasak/go-kubeadmclient/sshclient"
	"github.com/google/uuid"
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
		clientID:           uuid.New().String(),
	},
	}
}

func (n *MasterNode) ChangePermissionKubeconfig() error {
	return n.Run("sudo chown $USER:$USER /etc/kubernetes/admin.conf")
}

func (n *MasterNode) TaintAsMaster() error {
	return n.Run("KUBECONFIG=/etc/kubernetes/admin.conf kubectl taint nodes --all node-role.kubernetes.io/master-")
}

func (n *MasterNode) ApplyFile(file string) error {
	return n.Run("KUBECONFIG=/etc/kubernetes/admin.conf kubectl apply -f " + file)
}

func (n *MasterNode) ApplyFlannel() error {
	return n.Run("KUBECONFIG=/etc/kubernetes/admin.conf kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml")
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
		return "", err
	}

	return c["token"].(string), nil
}

func (n *MasterNode) Run(shell string) error {
	return n.sshClient().Run([]string{shell})
}

func (n *MasterNode) GetKubeConfig() (string, error) {
	return n.sshClient().Collect("sudo cat /etc/kubernetes/admin.conf")
}

func (n *MasterNode) GetJoinCommand() (string, error) {
	return n.sshClient().Collect("sudo kubeadm token create --print-join-command")
}

func (n *MasterNode) installAndFetchCommand(clusterName string, vip string) (string, error) {

	osType := n.determineOS()

	fmt.Println("os determined " + osType)

	var cmds []string

	if osType == Ubuntu {
		cmds = ubuntu.GenerateCommands()

		err := n.sshClient().Run(cmds)
		if err != nil {
			return "", err
		}

		err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(vip, clusterName)), "/tmp/kubeadm-config.yaml")
		if err != nil {
			return "", err
		}

	} else if osType == Centos || osType == RedHat {
		cmds = centos.GenerateCommands()
		err := n.sshClient().Run(cmds)
		if err != nil {
			return "", err
		}

		err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(n.ipOrHost, "")), "/tmp/kubeadm-config.yaml")
		if err != nil {
			return "", err
		}
	}

	out, err := n.sshClientWithTimeout(30 * time.Minute).Collect("sudo kubeadm init --config /tmp/kubeadm-config.yaml --upload-certs")
	if err != nil {
		log.Println(out)
		return "", err
	}

	return getControlPlaneJoinCommand(out), nil
}

func (n *MasterNode) Install(availability *common.HighAvailability) error {

	osType := n.determineOS()

	fmt.Println("os determined " + osType)

	var cmds []string

	if osType == Ubuntu {
		cmds = ubuntu.GenerateCommands()

		err := n.sshClientWithTimeout(30 * time.Minute).Run(cmds)
		if err != nil {
			return err
		}

	} else if osType == Centos || osType == RedHat {
		cmds = centos.GenerateCommands()
		err := n.sshClientWithTimeout(30 * time.Minute).Run(cmds)
		if err != nil {
			return err
		}
	}

	err := n.sshClientWithTimeout(30 * time.Minute).Run(cmds)
	if err != nil {
		return err
	}

	var s string

	if availability != nil {
		s = "sudo " + availability.JoinCommand
	} else {
		s = "sudo kubeadm init --pod-network-cidr=10.244.0.0/16"
	}

	return n.sshClientWithTimeout(30 * time.Minute).Run([]string{s})
}

func getControlPlaneJoinCommand(data string) string {
	var cmd string

	for _, line := range strings.Split(data, "\n") {

		if strings.HasPrefix(strings.TrimSpace(line), "kubeadm") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--discovery") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--control-plane") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
			return cmd
		}
	}

	return cmd
}
