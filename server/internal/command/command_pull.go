package command

import (
	"errors"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/provider"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"math"
	"strings"
	"sync"
)

func Pull(context *cli.Context) error {
	providerName := strings.TrimSpace(context.Args().Get(0))

	if providerName == "" {
		return errors.New("no provider name specified")
	}

	mediaProvider, err := provider.GetByName(providerName)

	if err != nil {
		return err
	}

	if mediaProvider == nil {
		return errors.New("media provider '" + providerName + "' not found.")
	}

	if err := mediaProvider.Init(); err != nil {
		return err
	}

	db, err := dbutils.GetConnection()

	if err != nil {
		return err
	}

	if err := dbutils.Migrate(db); err != nil {
		return err
	}

	mediaTypes := context.StringSlice(FlagPullType)

	mtMap := map[model.MediaType]bool{}

	for _, mt := range mediaTypes {
		mtMap[mt] = true
	}

	pages := int(math.Max(0, float64(context.Int(FlagPullPages))))

	wg := &sync.WaitGroup{}

	var lastError error

	for mt := range mtMap {
		wg.Add(1)

		go func(mediaType model.MediaType) {
			defer wg.Done()

			if err := mediaProvider.Pull(db, mediaType, pages); err != nil {
				log.Error(err)
				lastError = err
			}
		}(mt)
	}

	wg.Wait()

	return lastError
}
