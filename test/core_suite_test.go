package test_test

import (
	"fmt"
	"github.com/application-research/barge/core"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var app *cli.App

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

var _ = BeforeSuite(func() {
	//	load envi config
	godotenv.Load("./config/config.env")
	app = cli.NewApp()
	fmt.Println("BeforeSuite")
	app.Description = `'barge' is a cli tool to stream data to an existing Estuary node.`
	app.Name = "barge"
	app.Commands = []*cli.Command{
		core.InitCmd,
		core.ConfigCmd,
		core.LoginCmd,
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

	app.Run([]string{"barge", "init"})
	app.Run([]string{"barge", "config", "show"})
	//app.Run([]string{"barge", "config", "set", "token", os.Getenv("ESTUARY_API_KEY")})
	app.Run([]string{"barge", "login", os.Getenv("ESTUARY_API_KEY")})
	app.Run([]string{"barge", "config", "show"})

})

var _ = AfterSuite(func() {
	fmt.Println("AfterSuite")
	exec.Command("rm", "-rf", ".barge").Run()
	exec.Command("rm", "-rf", "~/.barge").Run()
})
