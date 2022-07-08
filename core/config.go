package core

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"
)

var ConfigCmd = &cli.Command{
	Name: "config",
	Subcommands: []*cli.Command{
		ConfigSetCmd,
		ConfigShowCmd,
	},
}

var ConfigSetCmd = &cli.Command{
	Name: "set",
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() != 2 {
			return fmt.Errorf("must pass two arguments: key and value")
		}

		viper.Set(cctx.Args().Get(0), cctx.Args().Get(1))

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
		return nil
	},
}

var ConfigShowCmd = &cli.Command{
	Name: "show",
	Action: func(cctx *cli.Context) error {
		var m map[string]interface{}
		if err := viper.Unmarshal(&m); err != nil {
			return err
		}

		b, err := json.MarshalIndent(m, "  ", "")
		if err != nil {
			return err
		}

		fmt.Println(string(b))
		return nil
	},
}
