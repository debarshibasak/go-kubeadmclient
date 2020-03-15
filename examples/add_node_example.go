// Package examples for adding nodes
package examples

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

// This is an example for adding node.
// In HA or no HA setup, you have to mention atleast one Master machine and the list of new nodes that have to be added
func AddNodeExample() {
	//Create clusters with only master machine
	k := kubeadmclient.Kubeadm{

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.51",
				"/Users//.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.55",
				"/Users//.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.56",
				"/Users//.ssh/id_rsa",
			),
		},
		Netorking:   networking.Flannel,
		VerboseMode: false,
	}

	err := k.AddNode()
	if err != nil {
		log.Fatal(err)
	}
}
