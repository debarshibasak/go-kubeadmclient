package kubeadmclient_test

import (
	"github.com/debarshibasak/kubekray/kubeadmclient"
	"log"
	"testing"
)

func TestKubeadm_CreateCluster(t *testing.T) {

	 k := kubeadmclient.Kubeadm{

	 	MasterNodes:[]*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode("ubuntu", "", ""),
			kubeadmclient.NewMasterNode("ubuntu", "", ""),
			kubeadmclient.NewMasterNode("ubuntu", "", ""),
		},

		 WorkerNodes:[]*kubeadmclient.WorkerNode{
			 kubeadmclient.NewWorkerNode("ubuntu", "", ""),
			 kubeadmclient.NewWorkerNode("ubun 	tu", "", ""),
			 kubeadmclient.NewWorkerNode("ubuntu", "", ""),
		 },

		 ApplyFiles:[]string{
			 "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		 },
	 }

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}