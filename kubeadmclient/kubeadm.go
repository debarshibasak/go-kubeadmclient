package kubeadmclient

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
	"github.com/pkg/errors"
)

//Reference - https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1

type Setup int

const (
	UNKNOWN Setup = 0
	HA      Setup = 1
	NONHA   Setup = 2
)

type Kubeadm struct {
	ClusterName       string
	MasterNodes       []*MasterNode
	WorkerNodes       []*WorkerNode
	HaProxyNode       *HaProxyNode
	ApplyFiles        []string
	PodNetwork        string
	ServiceNetwork    string
	DNSDomain         string
	VerboseMode       bool
	Netorking         *networking.Networking
	SkipWorkerFailure bool
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

func (k *Kubeadm) workerErrorManager(errc chan *workerError) error {
	for errWorker := range errc {
		if errWorker.err != nil {
			if errWorker.err == errWhileAddWorker {
				errWrk := errors.New("worker=" + errWorker.worker.ipOrHost + "err=" + errWorker.err.Error())
				if !k.SkipWorkerFailure {
					return errWrk
				}
				log.Println(errWrk.Error() + " however, skipping this error")
			} else {
				return errWorker.err
			}
		}
	}
	return nil
}
