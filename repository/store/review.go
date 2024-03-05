package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tabed23/travel-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewStore struct {
	review mongo.Collection
	tour   mongo.Collection
}

func NewReviewStore(r mongo.Collection, t mongo.Collection) *ReviewStore {
	return &ReviewStore{review: r, tour: t}
}

func (s *ReviewStore) CreateReviw(id string, review models.Review) (models.Review, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancle()
	var tour models.Tour
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Review{}, err
	}

	if err = s.tour.FindOne(ctx, bson.M{"_id": objId}).Decode(&tour); err != nil {
		return models.Review{}, err
	}
	review.ID = primitive.NewObjectID()
	review.TourID = tour.ID
	review.CreatedAt = time.Now().UTC()
	review.UpdatedAt = time.Now().UTC()
	arryFilter := bson.M{"$set": bson.M{"reviews": []interface{}{}}}
	if _, err := s.tour.UpdateOne(ctx, bson.M{"_id": objId}, arryFilter); err != nil {
		return models.Review{}, err
	}

	upsert := true
	filter := bson.M{"_id": objId}
	update := bson.M{
		"$push": bson.M{"reviews": review.ID},
	}
	opts := options.Update().SetUpsert(upsert)

	if _, err := s.tour.UpdateOne(ctx, filter, update, opts); err != nil {
		return models.Review{}, err
	}

	if _, err = s.review.InsertOne(ctx, review); err != nil {
		return models.Review{}, err
	}
	return review, nil
}

func (s *ReviewStore) IsExist(data, value string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	count, err := s.review.CountDocuments(ctx, bson.M{value: data})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (s *ReviewStore) GetAll(page, limit int) ([]models.Review, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	reviews := []models.Review{}

	cur, err := s.review.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []models.Review{}, 0, err
	}
	for cur.Next(ctx) {
		var review models.Review
		if err := cur.Decode(&review); err != nil {
			return []models.Review{}, 0, err
		}
		reviews = append(reviews, review)
	}
	if err := cur.Err(); err != nil {
		return []models.Review{}, 0, err
	}
	count, err := s.review.CountDocuments(ctx, bson.M{})
	return reviews, int(count), err
}

func (s *ReviewStore) CountReviews() (int, error) {
	opts := options.Count().SetHint("_id_")
	count, err := s.review.CountDocuments(context.TODO(), opts)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
func (s *ReviewStore) GetOne(id string) (models.Review, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Review{}, err
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var review models.Review
	if err := s.review.FindOne(ctx, bson.M{"_id": objId}).Decode(&review); err != nil {
		return models.Review{}, err
	}
	return review, nil
}

func (s *ReviewStore) Delete(id string) (bool, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.review.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return false, err
	}

	if deletedCount := result.DeletedCount; deletedCount == 0 {
		return false, errors.New("review not found")
	}
	if _, err := s.tour.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$pull": bson.M{"reviews": objId}}); err != nil {
		return false, fmt.Errorf("error removing review ID from tour: %v", err)
	}
	return true, nil

}
