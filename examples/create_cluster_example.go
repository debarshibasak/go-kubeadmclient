package examples

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

/*
This is an example for non-HA cluster creation.
This will create single controlplanes for the cluster and add the nodes to it.
*/
func CreateClusterExampleNonHA() {
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
		Netorking:   kubeadmclient.Flannel,
		VerboseMode: false,
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}

/*
This is an example for HA cluster creation.
This will create multiple controlplanes for the single cluster which can be accessed via an HAPRoxy server which is also provisioned.
*/
func CreateClusterExampleHA() {
	//Create clusters with only master machine
	k := kubeadmclient.Kubeadm{

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.51",
				"/Users//.ssh/id_rsa",
			),
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.50",
				"/Users//.ssh/id_rsa",
			),
			kubeadmclient.NewMasterNode("ubuntu",
				"192.168.64.52",
				"/Users//.ssh/id_rsa",
			),
		},
		HaProxyNode: kubeadmclient.NewHaProxyNode("ubuntu",
			"192.168.64.54",
			"/Users//.ssh/id_rsa",
		),
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
		Netorking:   kubeadmclient.Flannel,
		VerboseMode: false,
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}
