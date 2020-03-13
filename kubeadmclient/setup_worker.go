package kubeadmclient

import (
	"log"
	"sync"
)

func (k *Kubeadm) setupWorkers(joinCommand string) error {
	var workerWG sync.WaitGroup
	if len(k.WorkerNodes) > 0 {
		for _, workerNode := range k.WorkerNodes {

			workerWG.Add(1)

			go func(workerWG *sync.WaitGroup, node *WorkerNode) {
				if err := node.Install(joinCommand); err != nil {
					log.Println(err)
				}
				workerWG.Done()
			}(&workerWG, workerNode)
		}
	}

	workerWG.Wait()

	return nil
}
