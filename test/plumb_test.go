package test_test

import (
	"fmt"
	"github.com/application-research/barge/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
)

var _ = Describe("Plumb Tests", Ordered, func() {

	//	init
	plumbCmd := core.PlumbCmd
	app := cli.NewApp()

	BeforeAll(func() {
		fmt.Println("BeforeAll")
		app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
		app.Name = "barge"
		app.Commands = []*cli.Command{
			core.InitCmd,
			plumbCmd,
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

	//	basic config check
	It("check plumb name", func() {
		fmt.Println(plumbCmd)
		Expect(plumbCmd.Name).To(Equal("plumb"))
	})

	It("check plumb description", func() {
		Expect(plumbCmd.Description).To(Equal("low level plumbing commands"))
	})

	It("check plumb usage", func() {
		Expect(plumbCmd.Usage).To(Equal("plumb <command> [<args>]"))
	})

	It("check number of sub commands", func() {
		Expect(plumbCmd.Subcommands).To(HaveLen(4)) // increment if theres an additional subcommand
	})

	It("check plumb subcommand names", func() {
		var allSubCommandsThere = false
		for _, sub := range plumbCmd.Subcommands {
			if sub.Name == "put-file" || sub.Name == "put-dir" || sub.Name == "split-add" || sub.Name == "put-car" {
				allSubCommandsThere = true
			} else {
				Expect(allSubCommandsThere).To(BeTrue())
			}
		}
		Expect(allSubCommandsThere).To(BeTrue())
	})

	It("check plumb put-file", func() {
		Expect(plumbCmd.Subcommands[0].Name).To(Equal("put-file"))
		err := app.Run([]string{"barge", "plumb", "put-file", "files/put-file.text"})
		fmt.Println(err)
	})

	It("check plumb put-dir", func() {
		Expect(plumbCmd.Subcommands[1].Name).To(Equal("put-car"))
		err := app.Run([]string{"barge", "plumb", "put-car", "files/put-car.car"})
		fmt.Println(err)
	})

	It("check plumb put-car", func() {
		Expect(plumbCmd.Subcommands[2].Name).To(Equal("split-add"))
		err := app.Run([]string{"barge", "plumb", "split-add", "files/split-add.text"})
		fmt.Println(err)
	})

	It("check plumb split-add", func() {
		Expect(plumbCmd.Subcommands[3].Name).To(Equal("put-dir"))
		err := app.Run([]string{"barge", "plumb", "put-dir", "files/put-dir"})
		fmt.Println(err)
	})
})
