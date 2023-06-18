package command

import (
	"github.com/nitwhiz/movie-match/server/internal/sideload"
	"github.com/urfave/cli/v2"
)

func SideloadImport(_ *cli.Context) error {
	return sideload.Import()
}

func SideloadQuery(ctx *cli.Context) error {
	movieCount := ctx.Int("movie-count")
	tvCount := ctx.Int("tv-count")

	if movieCount == 0 {
		movieCount = 100
	}

	if tvCount == 0 {
		tvCount = 100
	}

	return sideload.Query(movieCount, tvCount)
}
