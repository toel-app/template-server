package db

import (
	"context"
	"fmt"
	"github.com/toel-app/registration/src/pkg/config"
	"github.com/toel-app/registration/src/pkg/logger"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	CollectionConversation = "conversations"
	CollectionMessage      = "messages"
)

const (
	OperatorAll      = "$all"
	OperatorSet      = "$set"
	OperatorNotEqual = "$ne"
	OperatorIn       = "$in"
)

var (
	once sync.Once
	db   *mongo.Database
)

func NewMongoDB(config config.Config, logger logger.Logger) *mongo.Database {
	once.Do(func() {
		var (
			username = config.Database.User
			password = config.Database.Pass
			host     = config.Database.Host
			database = config.Database.Name
			err      error
			uri      = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", username, password, host, database)
		)

		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Connect(ctx); err != nil {
			panic("Couldn't connect to mongo")
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			panic("Couldn't ping to mongo")
		}

		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		db = client.Database(database)
		logger.Info("mongodb initialized")
	})

	return db
}
