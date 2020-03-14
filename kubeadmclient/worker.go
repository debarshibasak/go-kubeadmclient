package kubeadmclient

import (
	"errors"
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

	if osType == nil {
		return errors.New("could not determine ostype, may be it could not ssh into it, or does not support the os")
	}

	if err := n.sshClientWithTimeout(30 * time.Minute).Run(osType.Commands()); err != nil {
		return err
	}

	return n.sshClientWithTimeout(30 * time.Minute).Run([]string{
		"sudo " + joinCommand,
	})
}
