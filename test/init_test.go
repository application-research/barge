package test_test

import (
	"fmt"
	"github.com/application-research/barge/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
)

var _ = Describe("Init Tests", Ordered, func() {
	//	init
	initCmd := core.InitCmd
	app := cli.NewApp()

	BeforeAll(func() {
		fmt.Println("BeforeAll")
		app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
		app.Name = "barge"
		app.Commands = []*cli.Command{
			initCmd,
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
		Expect(app.Run([]string{"barge", "help"})).To(Succeed())
	})

	It("init", func() {
		Expect(app.Run([]string{"barge", "init"})).To(Succeed())
	})
})
