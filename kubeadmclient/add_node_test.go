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
				"192.168.64.47",
				"USER_HOME/.ssh/id_rsa",
			),
		},
		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.51",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.52",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.50",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.48",
				"USER_HOME/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.49",
				"USER_HOME/.ssh/id_rsa",
			),
		},
		SkipWorkerFailure: true,
	}

	err := k.AddNode()
	if err != nil {
		t.Fatal(err)
	}
}
