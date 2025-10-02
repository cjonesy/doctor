package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cjonesy/doctor/internal/config"
	"github.com/cjonesy/doctor/internal/logger"
	"github.com/cjonesy/doctor/internal/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version: version.Version,
		Usage:   "checks your system for issues and suggests fixes",
		Action: func(cliCtx *cli.Context) error {
			// Create context with signal handling
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle signals for graceful shutdown
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sigChan
				fmt.Fprintln(os.Stderr, "\nReceived interrupt signal, shutting down...")
				cancel()
			}()

			// Create logger
			log := logger.New(
				cliCtx.Bool("verbose"),
				cliCtx.Bool("json"),
			)

			// Parse timeout
			timeout := cliCtx.Duration("timeout")
			if timeout > 0 {
				log.InfoContext(ctx, "timeout configured", "duration", timeout)
			}

			cfg := config.Config{
				Path:        cliCtx.Path("config"),
				Verbose:     cliCtx.Bool("verbose"),
				Timeout:     timeout,
				Parallelism: cliCtx.Int("parallel"),
				Logger:      log,
			}

			err := cfg.Run(ctx)
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
			&cli.DurationFlag{
				Name:    "timeout",
				Usage:   "timeout for each check (e.g., 30s, 1m)",
				Value:   0,
				Aliases: []string{"t"},
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output logs in JSON format",
			},
			&cli.IntFlag{
				Name:    "parallel",
				Usage:   "run checks in parallel (0=sequential, -1=unlimited, N=limit to N concurrent)",
				Value:   0, // Sequential by default
				Aliases: []string{"p"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
