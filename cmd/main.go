package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
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
			//TODO
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
