package kubeadmclient

import (
	"errors"
	"log"
	"sync"
	"time"
)

var (
	errNoWorkerForRemoveNode = errors.New("no worker information is set while removing node")
	errNoMasterForRemoveNode = errors.New("no master information is set while removing node")
)

// RemoveNode will take the incoming Kubeadm struct.
// For each worker, it just reset the configuration.
// By default it will fail if any of the worker fails to reset.
// However, you can skip that with the field SkipWorkerFailure in Kubeadm
func (k *Kubeadm) RemoveNode() error {
	startTime := time.Now()

	if len(k.WorkerNodes) == 0 {
		return errNoWorkerForRemoveNode
	}

	if len(k.MasterNodes) == 0 {
		return errNoMasterForRemoveNode
	}

	var hostnames []string
	var wg sync.WaitGroup

	for i, worker := range k.WorkerNodes {

		go func(worker *WorkerNode, i int, wg *sync.WaitGroup) {
			hostname, err := worker.drainAndReset()
			if err != nil {
				if !k.SkipWorkerFailure {
					log.Fatal(err)
				}
				log.Println(err)
			}

			if hostname != "" {
				hostnames = append(hostnames, hostname)
			}

		}(worker, i, &wg)
	}

	for _, hostname := range hostnames {
		if err := k.MasterNodes[0].deleteNode(hostname); err != nil {
			log.Println(err)
		}
	}

	log.Println("time taken = " + time.Since(startTime).String())

	return nil
}
