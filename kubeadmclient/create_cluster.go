package kubeadmclient

import "fmt"

type Kubeadm struct {
	MasterNodes []MasterNode
	WorkerNodes []WorkerNode
}

func (k *Kubeadm) CreateCluster() error {

	if len(k.MasterNodes) == 1 {
		//nonha setup

		masterNode := k.MasterNodes[0]
		err := masterNode.Install()
		if err != nil {
			return err
		}

		token, err := masterNode.GetToken()
		if err != nil {
			return err
		}

		fmt.Println(token)

		return nil
	}

	return nil
}
