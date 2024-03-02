package store

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	coll mongo.Collection
}

func NewUserStore(c mongo.Collection)*UserStore{
	return &UserStore{coll: c}
}