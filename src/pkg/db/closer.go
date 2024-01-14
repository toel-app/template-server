package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type Closer interface {
	Close() error
}

type closer struct {
	mongodb *mongo.Database
}

func NewCloser(mongodb *mongo.Database) Closer {
	return &closer{mongodb}
}

func (c closer) Close() error {
	if err := c.mongodb.Client().Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	fmt.Println("Database closed")
	return nil
}
