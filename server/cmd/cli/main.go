package main

import (
	"github.com/nitwhiz/movie-match/server/internal/command"
	"github.com/nitwhiz/movie-match/server/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "movie-match",
		Usage: "cli utility for the movie-match server",
		Commands: []*cli.Command{
			{
				Name:  "pull-media",
				Usage: "pull data from a media provider",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "media-type",
						Value: cli.NewStringSlice("all"),
						Usage: "which type of media to pull. possible values: movie, tv, all",
					},
				},
				ArgsUsage: "media-provider",
				Action:    command.PullMedia,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
