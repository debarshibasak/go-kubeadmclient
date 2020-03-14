package kubeadmclient_test

import (
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

func TestKubeadm_AddNode(t *testing.T) {
	//Create clusters with only master machine
	k := kubeadmclient.Kubeadm{
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
				"192.168.64.37",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.38",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		SkipAddWorkerFailure: true,
	}

	err := k.AddNode()
	if err != nil {
		t.Fatal(err)
	}
}
