package core

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var ConfigCmd = &cli.Command{
	Name:        "config",
	Description: `'barge config' is a command to set up the local barge configuration`,
	Usage:       "barge config <command>",
	Subcommands: []*cli.Command{
		ConfigSetCmd,
		ConfigShowCmd,
	},
}

var ConfigSetCmd = &cli.Command{
	Name:        "set",
	Description: `'barge config set <key> <value>' is a command to set up key value configuration'`,
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
	Name:        "show",
	Description: `'barge config show' is a command to show the existing configuration'`,
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
