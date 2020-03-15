//Package examples for deleting cluster
package examples

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

// This is an example for deleting cluster.
// This needs atleast one master to be specified
// This will fetch list of nodes and starting deleting them.
// Alternative you can set ResetOnDeleteCluster = true, this will also try to reset the node with best efforts.
func DeleteClusterExample() {
	k := kubeadmclient.Kubeadm{

		ClusterName: "testcluster",

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.2",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.3",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.10",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		SkipWorkerFailure:    false, //Skip if any of the worker has failure
		ResetOnDeleteCluster: true,  //Reset the node with best efforts if this field is set
		VerboseMode:          false, //Log level
	}

	err := k.DeleteCluster()
	if err != nil {
		log.Println(err)
	}
}
