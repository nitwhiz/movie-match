package command

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/urfave/cli/v2"
	"strings"
)

const FlagPullType = "type"
const FlagPullPages = "pages"

func GetApp() *cli.App {
	return &cli.App{
		Name:  "movie-match",
		Usage: "cli utility for the movie-match server",
		Commands: []*cli.Command{
			{
				Name:  "pull",
				Usage: "Pull media data from a provider",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  FlagPullType,
						Usage: "Which type of media to pull. Possible values: " + strings.Join(model.AvailableMediaTypes, ", "),
					},
					&cli.IntFlag{
						Name:  FlagPullPages,
						Value: 10,
						Usage: "How many pages to pull from the media provider, if it supports paging.",
					},
				},
				ArgsUsage: "media-provider",
				Action:    Pull,
			},
			{
				Name:   "serve",
				Usage:  "Start the server",
				Action: Server,
			},
		},
	}
}
