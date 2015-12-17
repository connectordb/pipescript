package main

import (
	"os"

	"github.com/connectordb/pipescript"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pipes"
	app.Usage = "Run the PipeScript data analysis engine on your given datasets."
	app.Version = pipescript.Version

	app.Copyright = "This software is available under the MIT license"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input,i",
			Usage: "The input file to perform analysis on.",
		},
		cli.StringFlag{
			Name:  "output,o",
			Value: "STDOUT",
			Usage: "File to write as output",
		},
	}

	app.Run(os.Args)

}