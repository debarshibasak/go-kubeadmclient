package kubeadmclient

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

//Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {

	var (
		joinCommand string
		err         error
	)

	if k.ClusterName == "" {
		return errors.New("cluster name is not set")
	}

	err = k.validateAndUpdateDefault()
	if err != nil {
		return err
	}

	startTime := time.Now()

	log.Println("total master - " + fmt.Sprintf("%v", len(k.MasterNodes)))
	log.Println("total workers - " + fmt.Sprintf("%v", len(k.WorkerNodes)))

	if k.HaProxyNode != nil {
		log.Println("total haproxy - " + fmt.Sprintf("%v", 1))
	}

	masterCreationStartTime := time.Now()
	joinCommand, err = k.setupMaster(k.determineSetup())
	if err != nil {
		return err
	}

	log.Printf("time taken to create masters = %v", time.Since(masterCreationStartTime))

	workerCreationTime := time.Now()

	if err := k.setupWorkers(joinCommand); err != nil {
		return err
	}

	log.Printf("time taken to create workers = %v", time.Since(workerCreationTime))

	for _, file := range k.ApplyFiles {
		err := k.MasterNodes[0].ApplyFile(file)
		if err != nil {
			return err
		}
	}

	if k.Netorking != nil {
		log.Printf("installing networking plugin = %v", k.Netorking.Name)
		err := k.MasterNodes[0].ApplyFile(k.Netorking.Manifests)
		if err != nil {
			return err
		}
	}

	log.Printf("Time taken to create cluster %v\n", time.Since(startTime).String())

	return nil
}
