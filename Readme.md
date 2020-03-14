# go-kubeadmclient

### What is go-kubeadmclient?

Golang SDK for creating kubernetes clusters. Operations that are currently supported are -

- Create Cluster
- Add nodes
- Remove nodes
 
The library automatically detects the operating system 
and then perform the operation for that particular os.

### Why is go-kubeadmclient?

Currently the only way to manage cluster operation on-prem/VM/baremetal machines are using custom ansible script or using kubespray.
`go-kubeadmclient` empowers you with a sdk to create cluster so that you can build your clusters based on logic/workflow etc. 
There are various projects for automation for custom kubernetes distribution. We believe something similar lacks for upstream kubernetes distribution. 
There aren't many SDK based approach to managing cluster operation therefore this project was initiated. 

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
	      Netorking:   networking.Flannel,
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

#### Removing node from existing cluster
if you want to remove node from an existing cluster follow this [example](https://github.com/debarshibasak/go-kubeadmclient/blob/master/examples/remove_node_example.go)
This example requires an existing cluster, access to master node and list of workers that has to be removed.

### Projects that use go-kubeadmclient
- [Kubestrike](https://github.com/debarshibasak/kubestrike)

#### Things to be noted
- This has been tested only on ubuntu. If you want this orchestration to be tested on centos or redhat, 
please create a github issue.
- Please make sure the user with which you create the VMs are passwordless sudoers.

#### Recent changes
- Added support for HA cluster
- Parallel worker node provisioning
- Added support for multipass
- Add HA Proxy Support for multi master setup
- cli for creating cluster
- More structure approach towards CNI, pod cidrs, service cidrs, dns domains etc.

#### Roadmap
- Testing this orchestration on centos, redhat
- Add Support for offline installation
- Add Support for remove node
- Add support for delete cluster
- Smart use MaxSession from ssh_config to buffer adding nodes

#### Supporting this project
- I need funding for testing this project
- If you want to join this project, please feel free to create pull requests or message me.
- You can support my effort with donation at [patreon](https://www.patreon.com/bePatron?u=31747625)

<a href="https://www.patreon.com/bePatron?u=31747625" data-patreon-widget-type="become-patron-button">Become a Patron!</a><script async src="https://c6.patreon.com/becomePatronButton.bundle.js"></script>