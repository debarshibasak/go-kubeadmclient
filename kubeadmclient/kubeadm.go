package kubeadmclient

import (
	"log"
	"sync"

	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
)

//Reference - https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1

type Setup int

const (
	UNKNOWN Setup = 0
	HA      Setup = 1
	NONHA   Setup = 2
)

type Kubeadm struct {
	ClusterName          string
	MasterNodes          []*MasterNode
	WorkerNodes          []*WorkerNode
	HaProxyNode          *HaProxyNode
	ApplyFiles           []string
	PodNetwork           string
	ServiceNetwork       string
	DNSDomain            string
	VerboseMode          bool
	Networking           *networking.Networking
	SkipWorkerFailure    bool
	ResetOnDeleteCluster bool
}

func (k *Kubeadm) GetKubeConfig() (string, error) {
	return k.MasterNodes[0].getKubeConfig()
}

func (k *Kubeadm) determineSetup() Setup {

	if len(k.MasterNodes) > 1 {
		return HA
	} else if len(k.MasterNodes) == 1 {
		return NONHA
	}

	return UNKNOWN
}

func (k *Kubeadm) deleteNodes(nodelist []string) error {

	var wg sync.WaitGroup
	if len(nodelist) > 0 {
		for _, node := range nodelist {
			wg.Add(1)
			go func(node string, wg *sync.WaitGroup) {
				if err := k.MasterNodes[0].deleteNode(node); err != nil {
					if !k.SkipWorkerFailure {
						log.Fatal(err)
					}

					log.Println(err)
				}

			}(node, &wg)
		}
	}
	wg.Wait()
	return nil
}

func (k *Kubeadm) validateAndUpdateDefault() error {

	if len(k.MasterNodes) == 0 {
		return errors.New("no master specified")
	}

	if k.ClusterName == "" {
		return errors.New("cluster name is empty")
	}

	if k.PodNetwork == "" {
		k.PodNetwork = "10.233.64.0/18"
	}

	if k.ServiceNetwork == "" {
		k.ServiceNetwork = "10.233.0.0/18"
	}

	if k.DNSDomain == "" {
		k.DNSDomain = "cluster.local"
	}

	return nil
}
