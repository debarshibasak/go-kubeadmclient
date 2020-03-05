package kubeadmclient

import (
	"encoding/json"
	"fmt"
	"github.com/debarshibasak/go-kubeadminclient/kubeadmclient/centos"
	"github.com/debarshibasak/go-kubeadminclient/kubeadmclient/common"
	"github.com/debarshibasak/go-kubeadminclient/kubeadmclient/ubuntu"
	"github.com/debarshibasak/go-kubeadminclient/sshclient"
	"github.com/google/uuid"
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
		clientID: uuid.New().String(),
	},
	}
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

func (n *MasterNode) GetKubeConfig() (string, error) {
	return n.sshClient().Collect("sudo cat /etc/kubernetes/admin.conf")
}

func (n *MasterNode) GetJoinCommand() (string, error) {
	return n.sshClient().Collect("sudo kubeadm token create --print-join-command")
}

func (n *MasterNode) InstallAndFetchCommand() (string, error) {

	osType := n.determineOS()

	fmt.Println("os determined " + osType)

	var cmds []string

	if osType == "ubuntu" {
		cmds = ubuntu.GenerateCommands(&common.HighAvailability{})

		err := n.sshClient().Run(cmds)
		if err != nil {
			return "", err
		}

		err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(n.ipOrHost)), "/tmp/kubeadm-config.yaml")
		if err != nil {
			return "", err
		}

	} else if osType == "centos" || osType == "redhat" {
		cmds = centos.GenerateCommands(&common.HighAvailability{})
		err := n.sshClient().Run(cmds)
		if err != nil {
			return "", err
		}

		err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(n.ipOrHost)), "/tmp/kubeadm-config.yaml")
		if err != nil {
			return "", err
		}
	}

	out, err := n.sshClientWithTimeout(20 * time.Minute).Collect("sudo kubeadm init --config /tmp/kubeadm-config.yaml --upload-certs")
	if err != nil {
		return "", err
	}

	return getControlPlaneJoinCommand(out), nil
}

func (n *MasterNode) Install(init bool, availability *common.HighAvailability) error {

	osType := n.determineOS()

	fmt.Println("os determined " + osType)

	var cmds []string

	if osType == "ubuntu" {
		cmds = ubuntu.GenerateCommands(availability)

		err := n.sshClientWithTimeout(30*time.Minute).Run(cmds)
		if err != nil {
			return err
		}

		if availability != nil && init {
			err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(n.ipOrHost)), "/tmp/kubeadm-config.yaml")
			if err != nil {
				return err
			}
		}

	} else if osType == "centos" || osType == "redhat" {
		cmds = centos.GenerateCommands(availability)
		err := n.sshClientWithTimeout(30*time.Minute).Run(cmds)
		if err != nil {
			return err
		}

		if availability != nil && init {
			err = n.sshClient().ScpToWithData([]byte(common.GenerateKubeadmConfig(n.ipOrHost)), "/tmp/kubeadm-config.yaml")
			if err != nil {
				return err
			}
		}
	}

	var s string

	if availability != nil && init {
		s = "sudo kubeadm init --config /tmp/kubeadm-config.yaml --upload-certs"
	} else if availability != nil && !init {
		s = "sudo " + availability.JoinCommand
	} else {
		s = "sudo kubeadm init --pod-network-cidr=192.178.0.0/16 --service-cidr=192.178.168.0/16"
	}

	return n.sshClientWithTimeout(30 * time.Minute).Run([]string{s})
}
