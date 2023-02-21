package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/internal/controller"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Option func(*Server, *gorm.DB) error

func WithAutoPull() Option {
	return func(server *Server, db *gorm.DB) error {
		autoPull, err := NewMediaAutoPull(db, "tmdb")

		if err != nil {
			log.Error("Auto Media Pull Init Error: ", err)
			return err
		}

		server.autoPull = autoPull

		return nil
	}
}

func WithTokenCleanup() Option {
	return func(server *Server, db *gorm.DB) error {
		server.tokenCleaner = auth.NewTokenCleanup(db)

		return nil
	}
}

func WithRouter() Option {
	return func(server *Server, db *gorm.DB) error {
		router, err := controller.Init(db)

		if err != nil {
			log.Error("Router Init Error: ", err)
			return err
		}

		server.router = router

		return nil
	}
}

type Server struct {
	router       *gin.Engine
	tokenCleaner *auth.TokenCleanup
	autoPull     *MediaAutoPull
}

func New(db *gorm.DB, opts ...Option) (*Server, error) {
	s := &Server{}

	for _, opt := range opts {
		if err := opt(s, db); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Start() error {
	if s.tokenCleaner != nil {
		s.tokenCleaner.Start()
	}

	if s.autoPull != nil {
		s.autoPull.Start()
	}

	if s.router != nil {
		addr := "0.0.0.0:6445"

		log.Infof("Listening on %s", addr)

		return s.router.Run(addr)
	}

	if s.tokenCleaner != nil {
		s.tokenCleaner.Wait()
	}

	if s.autoPull != nil {
		s.autoPull.Wait()
	}

	log.Info("bye!")

	return nil
}
