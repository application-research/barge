package tests_test

import (
	"bytes"
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
var outputBuffer bytes.Buffer

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Suite")
}

var _ = BeforeSuite(func() {
	//	load envi config
	godotenv.Load("./config/config.env")
	app = cli.NewApp()
	app.Writer = &outputBuffer
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
	app.After = func(cctx *cli.Context) error {
		fmt.Println("After1")
		fmt.Println(outputBuffer.String())
		fmt.Println("After2")
		return nil
	}
	app.Run([]string{"barge", "init"})
	app.Run([]string{"barge", "login", os.Getenv("ESTUARY_API_KEY")})

})

var _ = AfterSuite(func() {
	fmt.Println("AfterSuite")
	exec.Command("rm", "-rf", ".barge").Run()
	exec.Command("rm", "-rf", "~/.barge").Run()
})
