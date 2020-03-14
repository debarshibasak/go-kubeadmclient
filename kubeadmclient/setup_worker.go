package kubeadmclient

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

var errWhileAddWorker = errors.New("error while adding worker")

type WorkerError struct {
	worker *WorkerNode
	err    error
}

func (k *Kubeadm) setupWorkers(joinCommand string) error {
	var workerWG sync.WaitGroup
	errc := make(chan *WorkerError, 1)

	if len(k.WorkerNodes) > 0 {
		for i, workerNode := range k.WorkerNodes {

			workerWG.Add(1)

			go func(workerWG *sync.WaitGroup, node *WorkerNode, i int) {
				if err := node.Install(joinCommand); err != nil {
					errc <- &WorkerError{
						worker: node,
						err:    err,
					}
				}

				if i == len(k.WorkerNodes)-1 {
					close(errc)
				}
				workerWG.Done()
			}(&workerWG, workerNode, i)
		}
	}

	for errWorker := range errc {
		if errWorker.err != nil {
			if errWorker.err == errWhileAddWorker {
				errWrk := errors.New("worker=" + errWorker.worker.ipOrHost + "err=" + errWorker.err.Error())
				if !k.SkipAddWorkerFailure {
					return errWrk
				}
				log.Println(errWrk.Error() + " however, skipping this error")
			} else {
				return errWorker.err
			}
		}
	}

	workerWG.Wait()

	return nil
}
