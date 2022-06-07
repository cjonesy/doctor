package main

import (
	"fmt"
	"os"

	"github.com/cjonesy/doctor/internal/config"
	"github.com/cjonesy/doctor/internal/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version: version.Version,
		Usage:   "checks your system for issues and suggests fixes",
		Action: func(ctx *cli.Context) error {
			cfg := config.Config{
				Path:    ctx.Path("config"),
				Verbose: ctx.Bool("verbose"),
			}

			err := cfg.Run()
			return err
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "config",
				Usage:   "path to config file",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose output",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
