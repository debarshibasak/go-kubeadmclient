package kubeadmclient

import (
	"log"
	"testing"
)

func TestInstallMaster(t *testing.T) {
	masterNode := NewMasterNode("ubuntu", "192.168.64.2", "ubuntu", "")
	if err := masterNode.Install(); err != nil {
		log.Fatal(err)
	}

	token, err := masterNode.GetToken()
	if err != nil {
			log.Fatal(err)
	}

	_, err = masterNode.GetKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	workerNode := NewWorkerNode("ubuntu", "192.168.64.2", "ubuntu", "")
	workerNode.SetToken(token)

	if err := workerNode.Install(); err != nil {
		log.Fatal(err)
	}

}
