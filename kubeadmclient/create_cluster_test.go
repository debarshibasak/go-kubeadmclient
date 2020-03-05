package kubeadmclient_test

import (
	"github.com/debarshibasak/go-kubeadminclient/kubeadmclient"
	"log"
	"testing"
)

func TestKubeadm_CreateCluster3(t *testing.T) {
	//	data := `You can now join any number of the control-plane node running the following command on each as root:
	//
	//  kubeadm join 192.168.64.32:6443 --token hdov1n.953giwoejw47hyep \
	//    --discovery-token-ca-cert-hash sha256:efdc67e1c43c91d68aad4542a0e15af612638b0113ef851bf4c76523016dd1d8 \
	//    --control-plane --certificate-key e584418c55d4e7d6d4a4a81565968a8411e7d4c1efa2de7fb26bf8f3cf37ed68
	//
	//Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
	//As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use
	//"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.`
	//
	//
	//  fmt.Println(ge)

}
func TestKubeadm_CreateCluster2(t *testing.T) {

	//Create clusters with only master machine
	k := kubeadmclient.Kubeadm{

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.51",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.50",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewMasterNode("ubuntu",
				"192.168.64.52",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},
		ApplyFiles: []string{
			"https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		},
		VerboseMode: false,
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}

func TestKubeadm_CreateCluster(t *testing.T) {

	k := kubeadmclient.Kubeadm{

		MasterNodes: []*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode(
				"ubuntu",
				"192.168.64.18",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},

		WorkerNodes: []*kubeadmclient.WorkerNode{
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.18",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
			kubeadmclient.NewWorkerNode(
				"ubuntu",
				"192.168.64.18",
				"/Users/debarshibasak/.ssh/id_rsa",
			),
		},

		ApplyFiles: []string{
			"https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		},
	}

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}
