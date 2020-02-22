package kubeadmclient

import (
	"errors"
	"github.com/debarshibasak/go-kubeadminclient/kubectl"
	"log"
)

type Kubeadm struct {
	MasterNodes []*MasterNode
	WorkerNodes []*WorkerNode
	ApplyFiles []string
}

/*Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {

	var kubeCtl *kubectl.Kubectl
	var joinCommand string

	if len(k.MasterNodes) == 1 {
		//nonha setup

		masterNode := k.MasterNodes[0]
		err := masterNode.Install()
		if err != nil {
			return err
		}

		joinCommand, err = masterNode.GetJoinCommand()
		if err != nil {
			return err
		}

		kubeconfig, err := masterNode.GetKubeConfig()
		if err != nil {
			return err
		}

		kubeCtl = kubectl.New([]byte(kubeconfig))
		err = kubeCtl.TaintAllNodes("node-role.kubernetes.io/master-")
		if err != nil {
			return err
		}
	} else {
		return errors.New("not supported yet")
	}

	for _, workerNode := range k.WorkerNodes {
		if err := workerNode.Install(joinCommand); err != nil {
			log.Println(err)
		}
	}

	for _, file := range k.ApplyFiles {
		err := kubeCtl.ApplyFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}
