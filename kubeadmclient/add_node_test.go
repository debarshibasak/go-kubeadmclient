package kubeadmclient_test

import (
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

func TestKubeadm_AddNode(t *testing.T) {

	t.SkipNow()
	k := kubeadmclient.Kubeadm{
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
				"192.168.64.29",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		SkipWorkerFailure: true,
	}

	err := k.AddNode()
	if err != nil {
		t.Fatal(err)
	}
}
