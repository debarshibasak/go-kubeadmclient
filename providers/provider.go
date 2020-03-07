package providers

import (
	"github.com/debarshibasak/go-kubeadminclient/kubeadmclient"
	"github.com/debarshibasak/go-multipass/multipass"
	"github.com/pkg/errors"
)

func Get(provider string, mastercount int, workercount int) ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	var masters []string
	var workers []string

	var masterNodes []*kubeadmclient.MasterNode
	var workerNodes []*kubeadmclient.WorkerNode

	if provider == "multipass" {

		for i := 0; i < mastercount; i++ {
			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				return masterNodes, workerNodes, err
			}

			masters = append(masters, instance.IP)
		}

		for i := 0; i < workercount; i++ {
			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				return masterNodes, workerNodes, err
			}

			workers = append(workers, instance.IP)
		}

	} else {
		return masterNodes, workerNodes, errors.New("provider not supported")
	}

	for _, master := range masters {
		masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", master, "/home/user/.ssh/id_rsa"))
	}


	for _, worker := range workers {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", worker, "/home/user/.ssh/id_rsa"))
	}

	return masterNodes, workerNodes, nil
}
