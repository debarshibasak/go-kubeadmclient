package kubeadmclient_test

import (
	"log"
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

func TestKubeadm_CreateCluster2(t *testing.T) {

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
		Netorking:   kubeadmclient.Flannel,
		VerboseMode: false,
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}

func TestKubeadm_CreateCluster(t *testing.T) {

	k := kubeadmclient.Kubeadm{
		ClusterName: "test",
		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.33",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},

		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.34",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.35",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},

		SkipAddWorkerFailure: false,

		ApplyFiles: []string{
			"https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		},
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}
