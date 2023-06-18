package command

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/urfave/cli/v2"
	"strings"
)

const FlagPullType = "type"
const FlagPullPages = "pages"

const FlagServerWeb = "web"

const FlagServerWithTokenCleanup = "with-token-cleanup"

const FlagServerWithMediaAutoPull = "with-media-auto-pull"

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
				Name:   "purge",
				Usage:  "Remove all tables in the database",
				Action: Purge,
			},
			{
				Name:   "serve",
				Usage:  "Start the server",
				Action: Server,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  FlagServerWeb,
						Usage: "Enable web server",
					},
					&cli.BoolFlag{
						Name:  FlagServerWithTokenCleanup,
						Usage: "Enable automatic token cleanup on this server",
					},
					&cli.BoolFlag{
						Name:  FlagServerWithMediaAutoPull,
						Usage: "Enable automatic media pull on this server",
					},
				},
			},
			{
				Name:   "hash",
				Usage:  "Hash a string for usage as password in the user config",
				Action: Hash,
			},
			{
				Name:  "sideload",
				Usage: "Sideload all scraped media from mongodb and posters from filesystem. Intended to be used with tmdb-scraper",
				Subcommands: []*cli.Command{
					{
						Name:   "import",
						Usage:  "Import all media",
						Action: SideloadImport,
					},
					{
						Name:  "query",
						Usage: "Query for media",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name: "movie-count",
							},
							&cli.IntFlag{
								Name: "tv-count",
							},
						},
						Action: SideloadQuery,
					},
				},
			},
		},
	}
}
