package kubeadmclient

import (
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/common"
	"github.com/debarshibasak/go-kubeadmclient/kubectl"
	"log"
	"strings"
	"sync"
)

type Kubeadm struct {
	ClusterName string
	MasterNodes []*MasterNode
	WorkerNodes []*WorkerNode
	ApplyFiles []string
	PodNetwork string
	ServiceNetwork string
	VerboseMode bool
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

/*Creates cluster give a list of master nodes, worker nodes and then applies required kubernetes manifests*/
func (k *Kubeadm) CreateCluster() error {

	var kubeCtl *kubectl.Kubectl
	var joinCommand string

	primaryMaster := k.MasterNodes[0]
	if len(k.MasterNodes) == 1 {
		//nonha setup

		masterNode := primaryMaster
		masterNode.verboseMode = k.VerboseMode

		err := masterNode.Install(false,nil)
		if err != nil {
			return err
		}

		joinCommand, err = masterNode.GetJoinCommand()
		if err != nil {
			return err
		}

		kubeconfig, err := masterNode.GetKubeConfig()
		if err != nil {
			return err
		}

		kubeCtl = kubectl.New([]byte(kubeconfig))
		err = kubeCtl.TaintAllNodes("node-role.kubernetes.io/master-")
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

		kubeconfig, err := primaryMaster.GetKubeConfig()
		if err != nil {
			return err
		}

		for _, master := range k.MasterNodes[1:len(k.MasterNodes)] {
				err := master.Install(false, &common.HighAvailability{JoinCommand:masterJoinCommand})
				if err != nil {
					log.Println(err)
				}
		}

		kubeCtl = kubectl.New([]byte(kubeconfig))
		err = kubeCtl.TaintAllNodes("node-role.kubernetes.io/master-")
		if err != nil {
			return err
		}

		joinCommand, err = primaryMaster.GetJoinCommand()
		if err != nil {
			return err
		}
	}

	if len(k.WorkerNodes) > 0 {
		var workG sync.WaitGroup
		workG.Add(len(k.WorkerNodes))

		for _, workerNode := range k.WorkerNodes {
			go func(workG *sync.WaitGroup, workerNode WorkerNode) {
				if err := workerNode.Install(joinCommand); err != nil {
					log.Println(err)
				}
				workG.Done()
			}(&workG, *workerNode)
		}

		workG.Wait()
	}


	for _, file := range k.ApplyFiles {
		err := kubeCtl.ApplyFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}
