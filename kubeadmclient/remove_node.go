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

	var workerWG sync.WaitGroup
	errc := make(chan *workerError, 1)
	var hostnames []string

	for i, worker := range k.WorkerNodes {

		workerWG.Add(1)

		go func(wg *sync.WaitGroup, worker *WorkerNode, i int) {
			hostname, err := worker.drainAndReset()
			if err != nil {
				errc <- &workerError{
					worker: worker,
					err:    err,
				}
			}

			if hostname != "" {
				hostnames = append(hostnames, hostname)
			}

			if i == len(k.WorkerNodes)-1 {
				close(errc)
			}
			wg.Done()
		}(&workerWG, worker, i)
	}

	workerWG.Wait()

	e := k.workerErrorManager(errc)
	if e != nil {
		return e
	}

	for _, hostname := range hostnames {
		if err := k.MasterNodes[0].ctlCommand("kubectl delete node " + hostname); err != nil {
			log.Println(err)
		}
	}

	log.Println("time taken = " + time.Since(startTime).String())

	return nil
}
