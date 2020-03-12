package providers

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/common"
	"github.com/debarshibasak/go-multipass/multipass"
	"github.com/pkg/errors"
)

func Get(provider string, mastercount int, workercount int) ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error) {

	var masters []string
	var workers []string
	var haproxyIP string

	var masterNodes []*kubeadmclient.MasterNode
	var workerNodes []*kubeadmclient.WorkerNode
	var haproxy *kubeadmclient.HaProxyNode

	var publicKeyLocation string
	var privateKeyLocation string
	var err error

	if provider == "multipass" {

		publicKeyLocation, privateKeyLocation, err = common.PublicKeyExists()
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		publicKey, err := ioutil.ReadFile(publicKeyLocation)
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		if mastercount > 1 {
			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				return masterNodes, workerNodes, haproxy, err
			}

			err = multipass.Exec(&multipass.ExecRequest{
				Name:    instance.Name,
				Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
			})
			if err != nil {
				return masterNodes, workerNodes, haproxy, err
			}

			haproxyIP = instance.IP
		}

		for i := 0; i < mastercount; i++ {
			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				return masterNodes, workerNodes, haproxy, err
			}

			err = multipass.Exec(&multipass.ExecRequest{
				Name:    instance.Name,
				Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
			})
			if err != nil {
				return masterNodes, workerNodes, haproxy, err
			}

			masters = append(masters, instance.IP)
		}

		var workerWaitGroup sync.WaitGroup

		for i := 0; i < workercount; i++ {

			workerWaitGroup.Add(1)

			go func(workerWaitGroup *sync.WaitGroup) {
				defer workerWaitGroup.Done()

				instance, err := multipass.Launch(&multipass.LaunchReq{
					CPU: 2,
				})
				if err != nil {
					log.Println(err)
				}

				err = multipass.Exec(&multipass.ExecRequest{
					Name:    instance.Name,
					Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
				})
				if err != nil {
					log.Println(err)
				}

				workers = append(workers, instance.IP)

			}(&workerWaitGroup)
		}

		workerWaitGroup.Wait()

		for _, master := range masters {
			masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", master, privateKeyLocation))
		}

		if haproxyIP != "" {
			haproxy = kubeadmclient.NewHaProxyNode("ubuntu", haproxyIP, privateKeyLocation)
		}

		for _, worker := range workers {
			workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", worker, privateKeyLocation))
		}

	} else {
		return masterNodes, workerNodes, haproxy, errors.New("provider not supported")
	}

	return masterNodes, workerNodes, haproxy, nil
}
