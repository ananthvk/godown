package main

import (
	"context"
	"os"

	"github.com/ananthvk/godown/internal/download"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := (&cli.Command{
		Name:        "godown",
		Description: "godown is a concurrent file downloader",
		Version:     "0.0.1",
		ArgsUsage:   "<url>",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return cli.Exit("no urls specified", 1)
			}
			downloader := &download.Downloader{}
			downloader.Download(cmd.Args().First())
			return nil
		},
	})
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
