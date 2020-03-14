# go-kubeadmclient


### What is go-kubeadmclient?

Golang SDK for creating kubernetes clusters. You can add nodes a cluster too. The library automatically detects the operating system 
and then perform the operation for that particular os.

### Why is go-kubeadmclient?

Currently the only way to create clusters and add nodes on-prem/vms/baremetal machines is using ansible script or using kubespray.
`go-kubeadmclient` empowers you with sdk to create cluster so that you can build you logically build clusters.

#### Install
```
go get github.com/debarshibasak/go-kubeadmclient
```

#### SDK overview

#### Create a NonHA Cluster With code.

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
	      Netorking:   kubeadmclient.Flannel,
          VerboseMode: false,
	 }

	err := k.CreateCluster()
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Creating HA cluster
If you want to create and HA Cluster follow this [example](https://github.com/debarshibasak/go-kubeadmclient/blob/master/examples/create_cluster_example.go#L50)
The example sets up an HAProxy, master and workers.

#### Adding node to existing cluster
If you want to add nodes to an existing cluster follow this [example](https://github.com/debarshibasak/go-kubeadmclient/blob/master/examples/add_node_example.go)
The example requires an existing master that is setup and list of worker. The automation provisions workers and adds them to the cluster.

### Projects that use go-kubeadmclient
- [Kubestrike](https://github.com/debarshibasak/kubestrike)

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
- Support for offline installation


#### Supporting this project
- I need funding for testing this project
- If you want to join this project, please feel free to create pull requests.
- You can support my effort with donation at [patreon](https://www.patreon.com/bePatron?u=31747625)

<a href="https://www.patreon.com/bePatron?u=31747625" data-patreon-widget-type="become-patron-button">Become a Patron!</a><script async src="https://c6.patreon.com/becomePatronButton.bundle.js"></script>