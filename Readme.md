# go-kubeadmclient

Golang SDK for creating kubernetes clusters and kubectl. The library automatically detects the operating system 
and then creates the cluster for that particular os.


Create a NonHA Cluster With code. For example:

```
package main

import (
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"log"
)

func main(){

	 k := kubeadmclient.Kubeadm{

	 	MasterNodes:[]*kubeadmclient.MasterNode{
			kubeadmclient.NewMasterNode("ubuntu", "192.168.1.9", "/home/debarshi/.ssh/id_rsa"),
		},

		 WorkerNodes:[]*kubeadmclient.WorkerNode{
			 kubeadmclient.NewWorkerNode("ubuntu", "192.168.1.10", "/home/debarshi/.ssh/id_rsa"),
			 kubeadmclient.NewWorkerNode("ubuntu", "192.168.1.11", "/home/debarshi/.ssh/id_rsa"),
			 kubeadmclient.NewWorkerNode("ubuntu", "192.168.1.12", "/home/debarshi/.ssh/id_rsa"),
		 },

		 ApplyFiles:[]string{
			 "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		 },
	 }

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Using CLI to create clusters 

Build the project as follows

```
go build -o kubeorchestrator cmd/main.go 
./kubeorchestrator --provider multipass --master-count 2 --worker-count 2 --cluster-name test
```
This command will actually acquire 2+2+1 instances in multipass. 1 extra instance to provision HAProxy.
Currently the cli only supports multipass.

#### Breaking down the SDK

You can also create individual master and workers using the SDK too. For example, to create master node:
```
masterNode := kubeadmclient.NewMasterNode("ubuntu", "192.168.64.16", "/home/debarshibasak/.ssh/id_rsa")
if err := masterNode.Install(false); err != nil {
    log.Fatal(err)
}
```

To Fetch kubeconfig

```
kubeconfig, err := masterNode.GetKubeConfig()
if err != nil {
    log.Fatal(err)
}
```

Then taint the node as master

```
err = masterNode.TaintAsMaster()
if err != nil {
    log.Fatal(err)
}
```

Apply a CNI plugin

```
err = masterNode.ApplyFile("https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml")
if err != nil {
    log.Fatal(err)
}
```


Fetch the joining string from the master node

```
joinCommand, err := masterNode.GetJoinCommand()
if err != nil {
        log.Fatal(err)
}
```

..and then create worker nodes

```
workerNode := kubeadmclient.NewWorkerNode("ubuntu", "192.168.64.15", "/home/debarshibasak/.ssh/id_rsa")

if err := workerNode.Install(joinCommand); err != nil {
    log.Fatal(err)
}
```

#### Things to be noted
- This has been only tested on ubuntu. If you want this orchestration to be tested on centos or redhat, 
please create a github issue.
- Please make sure the user with which you create the VMs are passwordless sudoers.
- PodCidr is hardcoded. It will be remove and made into a variable in kubeadm very soon.

#### Recent changes
- Added support for HA cluster
- Parallel worker node provisioning
- Added support for multipass
- Add HA Proxy Support for multi master setup
- cli for creating cluster

#### Roadmap
- CLI support for baremetal, gke, aks, eks, digitalocean
- Testing this orchestration on centos
- use configuration file for the cli
- Support Multicloud providers, VM hypervisors
- More structure approach towards CNI, pod cidrs, service cidrs, etc.

