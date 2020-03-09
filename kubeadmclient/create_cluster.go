package kubeadmclient

import (
	"fmt"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/common"
	"log"
	"strings"
	"sync"
	"time"
)

type Kubeadm struct {
	ClusterName    string
	MasterNodes    []*MasterNode
	WorkerNodes    []*WorkerNode
	ApplyFiles     []string
	PodNetwork     string
	ServiceNetwork string
	VerboseMode    bool
}

func getControlPlaneJoinCommand(data string) string {
	var cmd string

	for _, line := range strings.Split(data, "\n") {

		if strings.HasPrefix(strings.TrimSpace(line), "kubeadm") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--discovery") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--control-plane") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
			return cmd
		}
	}

	return cmd
}

func (k *Kubeadm) GetKubeConfig() (string, error) {
	return k.MasterNodes[0].GetKubeConfig()
}

func (k *Kubeadm) ApplyTaint() (string, error) {
	return k.MasterNodes[0].GetKubeConfig()
}

//Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {

	var joinCommand string

	startTime := time.Now()

	log.Println("total master - " + fmt.Sprintf("%v", len(k.MasterNodes)))
	log.Println("total workers - " + fmt.Sprintf("%v", len(k.WorkerNodes)))

	primaryMaster := k.MasterNodes[0]
	if len(k.MasterNodes) == 1 {
		//nonha setup
		masterNode := primaryMaster
		masterNode.verboseMode = k.VerboseMode

		err := masterNode.Install(false, nil)
		if err != nil {
			return err
		}

		joinCommand, err = masterNode.GetJoinCommand()
		if err != nil {
			return err
		}

		err = masterNode.ChangePermissionKubeconfig()
		if err != nil {
			return err
		}

		err = masterNode.TaintAsMaster()
		if err != nil {
			return err
		}
	} else {

		primaryMaster.verboseMode = k.VerboseMode

		masterJoinCommand, err := primaryMaster.InstallAndFetchCommand()
		if err != nil {
			log.Println(err)
			return err
		}

		for _, master := range k.MasterNodes[1:len(k.MasterNodes)] {
			err := master.Install(false, &common.HighAvailability{JoinCommand: masterJoinCommand})
			if err != nil {
				log.Println(err)
			}
		}

		err = primaryMaster.ChangePermissionKubeconfig()
		if err != nil {
			return err
		}

		err = primaryMaster.TaintAsMaster()
		if err != nil {
			return err
		}

		joinCommand, err = primaryMaster.GetJoinCommand()
		if err != nil {
			return err
		}
	}

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

	for _, file := range k.ApplyFiles {

		err := k.MasterNodes[0].ApplyFile(file)
		if err != nil {
			return err
		}
	}

	log.Printf("Time taken to create cluster %v\n", time.Since(startTime).String())

	return nil
}
