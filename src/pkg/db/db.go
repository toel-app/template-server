package db

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/toel-app/common-utils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	COLLECTION_OTP = "otps"
)

var (
	Database *mongo.Database
	Context  context.Context
	Cancel   context.CancelFunc
	username = viper.GetString("database.user")
	password = viper.GetString("database.pass")
	host     = viper.GetString("database.host")
	database = viper.GetString("database.database_name")
)

func Connect() {
	var err error
	var uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", username, password, host, database)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	Context, Cancel = context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(Context)
	if err != nil {
		logger.Error("Couldn't connect to mongo", err)
		return
	}

	if err = client.Ping(Context, readpref.Primary()); err != nil {
		logger.Error("Couldn't ping to mongo", err)
		return
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	Database = client.Database(database)
	logger.Info("Mongodb OK")
}

func Close() {
	if Database == nil {
		return
	}

	err := Database.Client().Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to MongoDB closed")
}
