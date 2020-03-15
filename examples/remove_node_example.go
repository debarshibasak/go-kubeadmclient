// Package examples for removing nodes
package examples

import (
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
)

// This is an example for removing node from an existing cluster.
// You have to specify the one of the masters and the nodes you want to remove.
func RemoveNodeExample() {
	k := kubeadmclient.Kubeadm{
		ClusterName: "test",
		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.47",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.49",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.50",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.51",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},

		SkipWorkerFailure: false,
		Networking:        networking.Flannel,
	}

	if err := k.RemoveNode(); err != nil {
		t.Fatal(err)
	}
}
