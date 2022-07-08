package core

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"text/tabwriter"
)

var CollectionsCmd = &cli.Command{
	Name: "collections",
	Subcommands: []*cli.Command{
		CollectionsCreateCmd,
		CollectionsLsDirCmd,
	},
	Action: listCollections,
}

var CollectionsCreateCmd = &cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Required: true,
		},
		&cli.StringFlag{
			Name: "description",
		},
	},
	Action: func(cctx *cli.Context) error {
		c, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		col, err := c.CollectionsCreate(cctx.Context, cctx.String("name"), cctx.String("description"))
		if err != nil {
			return err
		}

		fmt.Println("new collection created")
		fmt.Println(col.Name)
		fmt.Println(col.UUID)

		return nil
	},
}

var CollectionsLsDirCmd = &cli.Command{
	Name:  "ls",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		c, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		if cctx.Args().Len() < 2 {
			return fmt.Errorf("must specify collection ID and path to list")
		}

		col := cctx.Args().Get(0)
		path := cctx.Args().Get(1)

		ents, err := c.CollectionsListDir(cctx.Context, col, path)
		if err != nil {
			return err
		}

		for _, e := range ents {
			if e.Dir {
				fmt.Println(e.Name + "/")
			} else {
				fmt.Println(e.Name)
			}
		}

		return nil
	},
}

func listCollections(cctx *cli.Context) error {
	c, err := LoadClient(cctx)
	if err != nil {
		return err
	}

	cols, err := c.CollectionsList(cctx.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	for _, c := range cols {
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\n", c.Name, c.UUID, c.CreatedAt)
		if err != nil {
			return err
		}
	}
	return w.Flush()
}
