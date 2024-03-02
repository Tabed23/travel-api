package store

import (
	"context"
	"time"

	"github.com/tabed23/travel-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewStore struct {
	coll mongo.Collection
}

func NewReviewStore(c mongo.Collection) *ReviewStore {
	return &ReviewStore{coll: c}
}

func (s *ReviewStore) CreateReviw(r models.Review) (models.Review, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	r.ID = primitive.NewObjectID()
	_, err := s.coll.InsertOne(ctx, r)
	if err != nil {
		return models.Review{}, err
	}
	return r, nil
}

func (s *ReviewStore) IsExist(data, value string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	count, err := s.coll.CountDocuments(ctx, bson.M{value: data})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
