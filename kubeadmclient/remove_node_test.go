package kubeadmclient_test

import (
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
)

func TestKubeadm_RemoveNode(t *testing.T) {

	k := kubeadmclient.Kubeadm{
		ClusterName: "test",
		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.47",
				"USER_HOME/.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.49",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.50",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.51",
				"USER_HOME/.ssh/id_rsa",
			),
		},

		SkipWorkerFailure: false,
		Netorking:         networking.Flannel,
	}

	if err := k.RemoveNode(); err != nil {
		t.Fatal(err)
	}
}
