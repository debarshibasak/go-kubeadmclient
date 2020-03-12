package kubeadmclient

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

type Setup int

const (
	UNKNOWN Setup = 0
	HA      Setup = 1
	NONHA   Setup = 2
)

type Kubeadm struct {
	ClusterName    string
	MasterNodes    []*MasterNode
	WorkerNodes    []*WorkerNode
	HaProxyNode    *HaProxyNode
	ApplyFiles     []string
	PodNetwork     string
	ServiceNetwork string
	VerboseMode    bool
}

func (k *Kubeadm) GetKubeConfig() (string, error) {
	return k.MasterNodes[0].GetKubeConfig()
}

func (k *Kubeadm) ApplyTaint() (string, error) {
	return k.MasterNodes[0].GetKubeConfig()
}

func (k *Kubeadm) determineSetup() Setup {

	if len(k.MasterNodes) > 1 {
		return HA
	} else if len(k.MasterNodes) == 1 {
		return NONHA
	}

	return UNKNOWN
}

//Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {
	if k.ClusterName == "" {
		return errors.New("cluster name is not set")
	}

	var (
		joinCommand string
		err         error
	)

	startTime := time.Now()

	log.Println("total master - " + fmt.Sprintf("%v", len(k.MasterNodes)))
	log.Println("total workers - " + fmt.Sprintf("%v", len(k.WorkerNodes)))

	if k.HaProxyNode != nil {
		log.Println("total haproxy - " + fmt.Sprintf("%v", 1))
	}

	masterCreationStartTime := time.Now()
	joinCommand, err = k.setupMaster(k.determineSetup())
	if err != nil {
		return err
	}

	log.Printf("time taken to create masters = %v", time.Since(masterCreationStartTime))

	workerCreationTime := time.Now()

	if err := k.setupWorkers(joinCommand); err != nil {
		return err
	}

	log.Printf("time taken to create workers = %v", time.Since(workerCreationTime))

	for _, file := range k.ApplyFiles {
		err := k.MasterNodes[0].ApplyFile(file)
		if err != nil {
			return err
		}
	}

	log.Printf("Time taken to create cluster %v\n", time.Since(startTime).String())

	return nil
}
