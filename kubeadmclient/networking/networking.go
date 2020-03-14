package networking

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
