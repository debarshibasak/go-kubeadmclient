package kubeadmclient_test

import (
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

func TestKubeadm_DeleteCluster(t *testing.T) {
	k := kubeadmclient.Kubeadm{

		ClusterName: "testcluster",

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.26",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.27",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.28",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		ResetOnDeleteCluster: true,
		SkipWorkerFailure:    true,
		VerboseMode:          false,
	}

	err := k.DeleteCluster()
	if err != nil {
		t.Fatal(err)
	}
}
