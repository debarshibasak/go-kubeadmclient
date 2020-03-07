# go-kubeadminclient

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


You can also create individual master and workers using the SDK too. For example, to create master node:

```
masterNode := NewMasterNode("ubuntu", "192.168.64.16", "/home/debarshibasak/.ssh/id_rsa")
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
kCtlClient := kubectl.New([]byte(kubeconfig))

err = kCtlClient.TaintAllNodes("node-role.kubernetes.io/master-")
if err != nil {
    log.Fatal(err)
}

```

Apply a CNI plugin

```
err = kCtlClient.ApplyFile("https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml")
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
workerNode := NewWorkerNode("ubuntu", "192.168.64.15", "/home/debarshibasak/.ssh/id_rsa")

if err := workerNode.Install(joinCommand); err != nil {
    log.Fatal(err)
}
```

#### Recent changes
- Added support for HA cluster
- Parallel worker node provisioning
- Added support for multipass

#### Roadmap
- cli for creating cluster
- use configuration file for the cli
- Support Multicloud providers, VM hypervisors

