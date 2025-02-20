package main

import (
	"fmt"
	"os"

	"github.com/ledgerwatch/erigon/common/dbutils"
	"github.com/ledgerwatch/erigon/log"
	"github.com/ledgerwatch/erigon/turbo/node"

	erigoncli "github.com/ledgerwatch/erigon/turbo/cli"

	"github.com/urfave/cli"
)

// defining a custom command-line flag, a string
var flag = cli.StringFlag{
	Name:  "custom-stage-greeting",
	Value: "default-value",
}

// defining a custom bucket name
const (
	customBucketName = "ch.torquem.demo.tgcustom.CUSTOM_BUCKET"
)

// the regular main function
func main() {
	// initializing Erigon application here and providing our custom flag
	app := erigoncli.MakeApp(runErigon,
		append(erigoncli.DefaultFlags, flag), // always use DefaultFlags, but add a new one in the end.
	)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Erigon main function
func runErigon(ctx *cli.Context) {
	// running a node and initializing a custom bucket with all default settings
	eri := node.New(ctx, node.Params{
		CustomBuckets: map[string]dbutils.BucketConfigItem{
			customBucketName: {},
		},
	})

	err := eri.Serve()

	if err != nil {
		log.Error("error while serving a Erigon node", "err", err)
	}
}
