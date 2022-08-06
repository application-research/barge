package main

import (
	"fmt"
	"github.com/application-research/barge/core"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
	app.Name = "barge"
	app.Commands = []*cli.Command{
		core.LoginCmd,
		core.InitCmd,
		core.ConfigCmd,
		core.BsGetCmd,
		core.PlumbCmd,
		core.CollectionsCmd,
		core.BargeAddCmd,
		core.BargeStatusCmd,
		core.BargeSyncCmd,
		core.BargeCheckCmd,
		core.BargeShareCmd,
		core.UiWebCmd,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logging",
		},
	}
	app.Before = func(cctx *cli.Context) error {
		if err := core.LoadConfig(); err != nil {
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
