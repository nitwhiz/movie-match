package server

import (
	"context"
	"github.com/nitwhiz/movie-match/server/internal/provider"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type MediaAutoPull struct {
	db       *gorm.DB
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	provider provider.MediaProvider
}

func NewMediaAutoPull(db *gorm.DB, providerName string) (*MediaAutoPull, error) {
	p, err := provider.GetMediaProviderByName(providerName)

	if err != nil {
		return nil, err
	}

	if err := p.Init(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MediaAutoPull{
		db:       db,
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
		provider: p,
	}, nil
}

func (p *MediaAutoPull) Pull() {
	log.Debug("Media Auto Pull Start")

	if err := p.provider.Pull(p.db, model.MediaTypeMovie, 50); err != nil {
		log.WithFields(log.Fields{
			"mediaType": model.MediaTypeMovie,
		}).Error("Auto Pull Error: ", err)
	}

	if err := p.provider.Pull(p.db, model.MediaTypeTv, 50); err != nil {
		log.WithFields(log.Fields{
			"mediaType": model.MediaTypeTv,
		}).Error("Auto Pull Error: ", err)
	}

	log.Debug("Media Auto Pull Finished")
}

func (p *MediaAutoPull) Start() {
	pullInterval := time.Hour * 6

	p.wg.Add(1)
	timer := time.NewTimer(pullInterval)

	go func() {
		defer p.wg.Done()
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				p.Pull()
				timer.Reset(pullInterval)
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

func (p *MediaAutoPull) Stop() {
	p.cancel()
	p.wg.Wait()
}

func (p *MediaAutoPull) Wait() {
	p.wg.Wait()
}
