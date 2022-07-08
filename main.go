package main

import (
	"fmt"
	"github.com/application-research/barge/core"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
	app.Name = "barge"
	app.Commands = []*cli.Command{
		core.ConfigCmd,
		core.LoginCmd,
		core.PlumbCmd,
		core.CollectionsCmd,
		core.InitCmd,
		core.BargeAddCmd,
		core.BargeStatusCmd,
		core.BargeSyncCmd,
		core.BargeCheckCmd,
		core.BargeShareCmd,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logging",
		},
	}
	app.Before = func(cctx *cli.Context) error {
		if err := loadConfig(); err != nil {
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() error {
	bargeDir, err := homedir.Expand("~/.barge")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(bargeDir, 0775); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.barge")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.WriteConfigAs(filepath.Join(bargeDir, "config"))
		} else {
			fmt.Printf("read err: %#v\n", err)
			return err
		}
	}
	return nil
}
