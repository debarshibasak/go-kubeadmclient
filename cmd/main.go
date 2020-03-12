package main

import (
	"fmt"
	"log"
	"os"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/providers"
	"github.com/urfave/cli"
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
			&cli.StringFlag{
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
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose mode",
			},
		},
		Name:  "create",
		Usage: "run the manifest",
		Action: func(c *cli.Context) error {
			provider := c.String("provider")
			fmt.Println(provider)

			log.Println("creating vm...")

			provider = c.String("provider")

			if provider == "" {
				log.Fatal("provider is not set")
			}

			masterNodes, workerNodes, haproxy, err := providers.Get(
				provider,
				c.Int("master-count"),
				c.Int("worker-count"),
			)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("creating cluster...")

			kubeadmClient := kubeadmclient.Kubeadm{
				ClusterName: c.String("cluster-name"),
				HaProxyNode: haproxy,
				MasterNodes: masterNodes,
				WorkerNodes: workerNodes,
				ApplyFiles: []string{
					"https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
				},
				VerboseMode: c.Bool("verbose"),
			}

			err = kubeadmClient.CreateCluster()
			if err != nil {
				log.Fatal(err)
			}

			printSummary(kubeadmClient)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printSummary(kubeadm kubeadmclient.Kubeadm) {
	fmt.Println("master machines")
	fmt.Println("-----------")
	for _, master := range kubeadm.MasterNodes {
		fmt.Println(master)
	}

	if kubeadm.HaProxyNode != nil {
		fmt.Println("-----------")
		fmt.Println("haproxy machines")
		fmt.Println("-----------")
		fmt.Println(kubeadm.HaProxyNode)
	}

	fmt.Println("-----------")
	fmt.Println("workers machines")
	fmt.Println("-----------")

	for _, worker := range kubeadm.WorkerNodes {
		fmt.Println(worker)
	}

}
