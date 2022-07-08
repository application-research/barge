package core

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"net/url"
)

var LoginCmd = &cli.Command{
	Name: "login",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "https://api.estuary.tech",
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() {
			return fmt.Errorf("must specify api token")
		}

		tok := cctx.Args().First()

		ec := &EstClient{
			Host: cctx.String("host"),
			Tok:  tok,
		}

		vresp, err := ec.Viewer(cctx.Context)
		if err != nil {
			return err
		}

		fmt.Println("logging in as user: ", vresp.Username)

		if len(vresp.Settings.UploadEndpoints) > 0 {
			sh := vresp.Settings.UploadEndpoints[0]
			u, err := url.Parse(sh)
			if err != nil {
				return err
			}

			u.Path = ""
			u.RawQuery = ""
			u.Fragment = ""

			fmt.Printf("selecting %s as our primary shuttle\n", u.String())

			viper.Set("estuary.primaryShuttle", u.String())
		}

		viper.Set("estuary.token", tok)
		viper.Set("estuary.host", ec.Host)

		return viper.WriteConfig()
	},
}
