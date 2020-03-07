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
		},
		Name:  "create",
		Usage: "run the manifest",
		Action: func(c *cli.Context) error {
			provider := c.String("provider")
			fmt.Println(provider)

			masterNodes, workerNodes,  err := providers.Get(c.String("provider"), c.Int("master-count"), c.Int("worker-count"))
			if err != nil {
				log.Fatal(err)
			}

			kubeadmClient := kubeadmclient.Kubeadm{
				ClusterName: c.String("cluster-name"),
				MasterNodes: masterNodes,
				WorkerNodes: workerNodes,
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
