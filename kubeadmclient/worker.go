package kubeadmclient

import (
	"time"

	"github.com/google/uuid"
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
			clientID:           uuid.New().String(),
		},
	}
}

func (n *WorkerNode) Install(joinCommand string) error {

	osType := n.determineOS()

	if err := n.sshClientWithTimeout(30 * time.Minute).Run(osType.Commands()); err != nil {
		return err
	}

	return n.sshClientWithTimeout(30 * time.Minute).Run([]string{
		"sudo " + joinCommand,
	})
}
