package store

import (
	"context"
	"fmt"
	"log/slog"
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
	logger *slog.Logger
}

func NewReviewStore(r mongo.Collection, t mongo.Collection, l *slog.Logger) *ReviewStore {
	return &ReviewStore{review: r, tour: t, logger: l}
}

// Create a new review Document
func (s *ReviewStore) CreateReviw(id string, review models.Review) (models.Review, error) {
	s.logger.Info("repository", "CreateReviw", "Review")

	ctx, cancle := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancle()
	var tour models.Tour
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("repository", "CreateReviw", err.Error())
		return models.Review{}, err
	}

	if err = s.tour.FindOne(ctx, bson.M{"_id": objId}).Decode(&tour); err != nil {
		s.logger.Error("repository", "CreateReviw", err.Error())
		return models.Review{}, err
	}

	s.logger.Info("repository", "CreateReviw", fmt.Sprintf("Document Found %v", tour))

	review.ID = primitive.NewObjectID()
	review.TourID = tour.ID
	review.CreatedAt = time.Now().UTC()
	review.UpdatedAt = time.Now().UTC()
	arryFilter := bson.M{"$set": bson.M{"reviews": []interface{}{}}}
	if _, err := s.tour.UpdateOne(ctx, bson.M{"_id": objId}, arryFilter); err != nil {
		s.logger.Error("repository", "CreateReviw", err.Error())
		return models.Review{}, err
	}

	upsert := true
	filter := bson.M{"_id": objId}
	update := bson.M{
		"$push": bson.M{"reviews": review.ID},
	}
	opts := options.Update().SetUpsert(upsert)

	if _, err := s.tour.UpdateOne(ctx, filter, update, opts); err != nil {
		s.logger.Error("repository", "CreateReviw", err.Error())
		return models.Review{}, err
	}
	s.logger.Info("repository", "CreateReviw", fmt.Sprintf("Document Review insetred in Tour %v", review))

	if _, err = s.review.InsertOne(ctx, review); err != nil {
		s.logger.Error("repository", "CreateReviw", err.Error())
		return models.Review{}, err
	}
	s.logger.Info("repository", "CreateReviw", fmt.Sprintf("Document Review Created %v", review))

	return review, nil
}

// ISExist Document
func (s *ReviewStore) IsExist(data, value string) (bool, error) {
	s.logger.Info("repository", "IsExist", "Review")

	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	count, err := s.review.CountDocuments(ctx, bson.M{value: data})
	if err != nil {
		s.logger.Error("repository", "IsExist", err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	s.logger.Warn("repository", "IsExist", fmt.Sprintf("%v documents not found", value))

	return false, nil
}

// Get All Document
func (s *ReviewStore) GetAll(page, limit int) ([]models.Review, int, error) {
	s.logger.Info("repository", "GetAll", "Review")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	reviews := []models.Review{}

	cur, err := s.review.Find(ctx, bson.M{}, opts)
	if err != nil {
		s.logger.Error("repository", "GetAll", err.Error())
		return []models.Review{}, 0, err
	}
	for cur.Next(ctx) {
		var review models.Review
		if err := cur.Decode(&review); err != nil {
			s.logger.Error("repository", "GetAll", err.Error())

			return []models.Review{}, 0, err
		}
		reviews = append(reviews, review)
	}
	if err := cur.Err(); err != nil {
		s.logger.Error("repository", "GetAll", err.Error())
		return []models.Review{}, 0, err
	}
	count, err := s.review.CountDocuments(ctx, bson.M{})
	return reviews, int(count), err
}

// Count Reviews Documents
func (s *ReviewStore) CountReviews() (int, error) {
	s.logger.Info("repository", "CountReviews", "Review")

	opts := options.Count().SetHint("_id_")
	count, err := s.review.CountDocuments(context.TODO(), opts)
	if err != nil {
		s.logger.Error("repository", "CountReviews", err.Error())
		return 0, err
	}
	s.logger.Error("repository", "CountReviews", fmt.Sprintf("Found %d Documents", count))

	return int(count), nil
}

// Get One Review Document
func (s *ReviewStore) GetOne(id string) (models.Review, error) {
	s.logger.Info("repository", "GetOne", "Review")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Review{}, err
	}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var review models.Review
	if err := s.review.FindOne(ctx, bson.M{"_id": objId}).Decode(&review); err != nil {
		s.logger.Error("repository", "GetOne", err.Error())
		return models.Review{}, err
	}
	s.logger.Info("repository", "GetOne", fmt.Sprintf("Found %v", review))

	return review, nil
}

// Delete Document
func (s *ReviewStore) Delete(tourId, reviewId string) (bool, error) {
	s.logger.Info("repository", "Delete", "Review")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tourObjId, err := primitive.ObjectIDFromHex(tourId)
	if err != nil {
		return false, err
	}
	reviewObjId, err := primitive.ObjectIDFromHex(reviewId)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": tourObjId}
	update := bson.M{"$pull": bson.M{"reviews": reviewObjId}}
	if _, err := s.tour.UpdateOne(ctx, filter, update); err != nil {
		s.logger.Error("repository", "Delete", err.Error())
		return false, err
	}

	if _, err := s.review.DeleteOne(ctx, bson.M{"_id": reviewObjId}); err != nil {
		s.logger.Error("repository", "Delete", err.Error())
		return false, err
	}
	s.logger.Info("repository", "Delete", fmt.Sprintf("reviews: %v deleted", reviewObjId))

	return true, nil
}
