package auth

import (
	"context"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

type TokenCleanup struct {
	db     *gorm.DB
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewTokenCleanup(db *gorm.DB) *TokenCleanup {
	ctx, cancel := context.WithCancel(context.Background())

	return &TokenCleanup{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

func (c *TokenCleanup) Purge() {
	if err := c.db.Where("valid_until <= NOW()").Delete(&model.UserToken{}).Error; err != nil {
		log.Error("Token Cleanup Error: " + err.Error())
	}
}

func (c *TokenCleanup) Start() {
	c.wg.Add(1)
	timer := time.NewTicker(time.Second * 10)

	go func() {
		defer c.wg.Done()
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				c.Purge()
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *TokenCleanup) Stop() {
	c.cancel()
	c.wg.Wait()
}
