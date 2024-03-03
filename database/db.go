package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	db *mongo.Client
}

func NewDatabase(monguri string) (*MongoDB, error) {

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(monguri))
	if err != nil {
		log.Fatal("error connecting to MongoDB")
		return nil, err
	}
	fmt.Println("Connecting to MongoDB")
	return &MongoDB{db}, nil
}

func (m *MongoDB) Close() {
	m.db.Disconnect(context.TODO())
}

func (m *MongoDB) GetDB() *mongo.Database {
	return m.db.Database("travle")
}
