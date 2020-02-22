# go-kubeadminclient

Golang sdk for creating kubernetes clusters and kubectl. The library automatically detects the operating system 
and then creates the cluster for that particular os.

For creating master node:

```
 	masterNode := NewMasterNode("ubuntu", "192.168.64.16", "/home/debarshibasak/.ssh/id_rsa")
 	if err := masterNode.Install(); err != nil {
 		log.Fatal(err)
 	}

```

Fetch kubeconfig

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

#### Roadmap

- Add High availability cluster creation
- Parallel worker node provisioning
- cli for creating cluster
- use configuration file for the cli

