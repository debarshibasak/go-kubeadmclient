package kubeadmclient

import (
	"log"
	"sync"

	"errors"
)

var errWhileAddWorker = errors.New("error while adding worker")

func (k *Kubeadm) setupWorkers(joinCommand string) error {

	var wg sync.WaitGroup
	if len(k.WorkerNodes) > 0 {
		for i, workerNode := range k.WorkerNodes {
			wg.Add(1)
			go func(node *WorkerNode, i int, wg *sync.WaitGroup) {
				defer wg.Done()
				err := node.install(joinCommand)
				if err != nil {
					if !k.SkipWorkerFailure {
						log.Fatal(err)
					}
					log.Println(err)
				}
			}(workerNode, i, &wg)
		}
	}

	wg.Wait()

	return nil
}
