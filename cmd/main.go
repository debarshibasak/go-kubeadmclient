package main

import (
	"fmt"
	"github.com/debarshibasak/kubekray/kubeadmclient"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ostype",
				Usage: "location of the script",
			},
			&cli.StringFlag{
				Name:  "username",
				Usage: "location of the script",
			},
		},
		Name:  "run",
		Usage: "run the manifest",
		Action: func(c *cli.Context) error {

			ostype := c.String("ostype")
			username := c.String("username")

			fmt.Println(ostype)
			d, err := ioutil.ReadFile("host.env")
			if err != nil {
				log.Fatal(err)
			}

			hosts := strings.Split(string(d), "\n")

			masterNode := kubeadmclient.NewMasterNode(username, hosts[0], ostype, "")
			if err := masterNode.Install(); err != nil {
				log.Println("error while installing master")
				log.Fatal(err)
			}

			//kubeConfig , err := masterNode.GetKubeConfig()
			//if err := masterNode.Install(); err != nil {
			//	log.Println("error while installing master")
			//	log.Fatal(err)
			//}

			token, err := masterNode.GetToken()
			if err := masterNode.Install(); err != nil {
				log.Println("error while installing master")
				log.Fatal(err)
			}

			for _, host := range hosts[1:] {
				workerNode := kubeadmclient.NewWorkerNode(username, host, ostype, "")
				workerNode.SetToken(token)
				if err := workerNode.Install(); err != nil {
					log.Println("error while installing master")
					log.Fatal(err)
				}
			}

			//kubeadmclient.CreatePodNetwork(kubeConfig)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
