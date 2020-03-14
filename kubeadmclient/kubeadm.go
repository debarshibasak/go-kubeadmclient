package kubeadmclient

import "github.com/pkg/errors"

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
