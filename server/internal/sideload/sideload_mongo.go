package sideload

import (
	"context"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func getMongoConnection() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	return mongo.Connect(ctx, options.Client().ApplyURI(config.C.Sideload.MongoDb.Uri))
}
