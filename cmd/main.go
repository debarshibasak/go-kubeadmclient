package main

import (
	"fmt"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/providers"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "provider",
				Usage: "set a provider",
			},
			&cli.IntFlag{
				Name:  "master-count",
				Usage: "master count",
			},
			&cli.IntFlag{
				Name:  "worker-count",
				Usage: "worker count",
			},
			&cli.IntFlag{
				Name:  "cluster-name",
				Usage: "name of the cluster",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "configuration file",
			},
			&cli.StringFlag{
				Name:  "cni",
				Usage: "choose the networking layer",
			},
		},
		Name:  "create",
		Usage: "run the manifest",
		Action: func(c *cli.Context) error {
			provider := c.String("provider")
			fmt.Println(provider)

			//TODO add preflight check
			//check if kubectl exists
			//check if multipass exists

			log.Println("creating vm...")

			//Check the CNI, default cni is flannel, it also has an effect on pod cidr
			//var cni = c.String("cni")
			//if cni == "" {
			//	cni = "flannel"
			//}

			masterNodes, workerNodes, err := providers.Get(c.String("provider"), c.Int("master-count"), c.Int("worker-count"))
			if err != nil {
				log.Fatal(err)
			}

			log.Println("creating cluster...")

			kubeadmClient := kubeadmclient.Kubeadm{
				ClusterName: c.String("cluster-name"),
				MasterNodes: masterNodes,
				WorkerNodes: workerNodes,
				ApplyFiles: []string{
					"https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
				},
				VerboseMode: false,
			}

			err = kubeadmClient.CreateCluster()
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
