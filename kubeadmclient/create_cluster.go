package kubeadmclient

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

//Reference - https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1

type Setup int

const (
	UNKNOWN Setup = 0
	HA      Setup = 1
	NONHA   Setup = 2
)

type Networking struct {
	Manifests string
	Name      string
}

func LookupNetworking(cni string) *Networking {
	switch cni {
	case "flannel":
		return Flannel
	case "canal":
		return Canal
	case "Calico":
		return Calico
	default:
		return nil
	}
}

var (
	Flannel = &Networking{
		Manifests: "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		Name:      "flannel",
	}

	Canal = &Networking{
		Manifests: "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		Name:      "canal",
	}

	Calico = &Networking{
		Manifests: "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		Name:      "calico",
	}
)

type Kubeadm struct {
	ClusterName    string
	MasterNodes    []*MasterNode
	WorkerNodes    []*WorkerNode
	HaProxyNode    *HaProxyNode
	ApplyFiles     []string
	PodNetwork     string
	ServiceNetwork string
	DNSDomain      string
	VerboseMode    bool
	Netorking      *Networking
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

//Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {

	var (
		joinCommand string
		err         error
	)

	if k.ClusterName == "" {
		return errors.New("cluster name is not set")
	}

	err = k.validateAndUpdateDefault()
	if err != nil {
		return err
	}

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

	if k.Netorking != nil {
		log.Printf("installing networking plugin = %v", k.Netorking.Name)
		err := k.MasterNodes[0].ApplyFile(k.Netorking.Manifests)
		if err != nil {
			return err
		}
	}

	log.Printf("Time taken to create cluster %v\n", time.Since(startTime).String())

	return nil
}
