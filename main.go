package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ananthvk/godown/internal/download"
	"github.com/urfave/cli/v3"
	// "github.com/vbauerster/mpb/v8"
	// "github.com/vbauerster/mpb/v8/decor"
)

func main() {
	cmd := (&cli.Command{
		Name:        "godown",
		Description: "godown is a concurrent file downloader",
		Version:     "0.0.1",
		ArgsUsage:   "<url>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output-dir",
				Value: ".",
				Usage: "directory to save files to",
			},
			&cli.BoolFlag{
				Name:  "ignore-invalid-url",
				Value: false,
				Usage: "ignores invalid urls that are passed as input, if the input url is missing a scheme, automatically prepends http://",
			},
			&cli.BoolFlag{
				Name:  "log",
				Value: false,
				Usage: "Enables logging",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return cli.Exit("no urls specified", 1)
			}

			if !cmd.Bool("log") {
				slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
			}

			downloader := download.NewDownloader(cmd.String("output-dir"), cmd.Bool("ignore-invalid-url"))

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer cancel()

			for _, url := range cmd.Args().Slice() {
				downloader.Download(ctx, url)
			}

			slog.Info("waiting for all downloads to complete")
			downloader.Wait()
			slog.Info("completed all downloads")
			return nil
		},
	})
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
