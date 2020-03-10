package common

import (
	"os"
)

type HighAvailability struct {
	JoinCommand string
}

func PublicKeyExists() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	publicKeyLocation := homeDir + "/.ssh/id_rsa.pub"
	privateKeyLocation := homeDir + "/.ssh/id_rsa"
	if _, err := os.Stat(publicKeyLocation); err == nil {
		if _, err := os.Stat(privateKeyLocation); err == nil {
			return publicKeyLocation, privateKeyLocation, nil
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}
}

func GenerateKubeadmConfig(ip string, clusterName string) string {
	return `
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: "1.17.3"
apiServer:
   certSANs:
   - "` + ip + `"
controlPlaneEndpoint: "` + ip + `:6443"
networking:
  podSubnet: 10.244.0.0/16
  clusterName: "` + clusterName + `"
`
}

//https://www.jordyverbeek.nl/nieuws/kubernetes-ha-cluster-installation-guide
func GenerateKubeInitConfig(ip string) string {

	return `
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: stable
apiServerCertSANs:
- "k8s.apiserver.cluster"
- ` + ip + `
controlPlaneEndpoint: "k8s.apiserver.cluster:443"
etcd:
  local:
    extraArgs:
      listen-client-urls: "https://127.0.0.1:2379,https://` + ip + `:2379"
      advertise-client-urls: "https://` + ip + `:2379"
      listen-peer-urls: "https://` + ip + `:2380"
      initial-advertise-peer-urls: "https://` + ip + `:2380"
      initial-cluster: "k8s-master01=https://` + ip + `:2380"
    serverCertSANs:
      - ` + ip + `
    peerCertSANs:
      - ` + ip + `
apiServerExtraArgs:
  service-node-port-range: 8000-31274
networking:
    # For flannel use 10.244.0.0/16, calico uses 192.168.0.0/16
    podSubnet: "10.244.0.0/16"
`
}
