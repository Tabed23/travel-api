package store

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tabed23/travel-api/models"
	"github.com/tabed23/travel-api/utils/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TourStore struct {
	coll   mongo.Collection
	logger *slog.Logger
}

func NewTourStore(c mongo.Collection, l *slog.Logger) *TourStore {
	return &TourStore{coll: c, logger: l}
}
func (t *TourStore) CreateTour(tour models.Tour) (models.Tour, error) {
	t.logger.Info("repository", "CreateTour", "Tour")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tour.ID = primitive.NewObjectID()
	tour.CreatedAt = time.Now().UTC()
	tour.UpdatedAt = time.Now().UTC()

	_, err := t.coll.InsertOne(ctx, tour)
	if err != nil {
		return models.Tour{}, err
	}

	return tour, nil
}
func (t *TourStore) IsExist(ctx context.Context, filter primitive.M) (bool, error) {
	t.logger.Info("repository", "IsExist", "Tour")

	count, err := t.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}

func (t *TourStore) DeleteTour(_id string) (bool, error) {
	t.logger.Info("repository", "DeleteTour", "Tour")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(_id)
	filter := bson.M{"_id": objId}
	ok, err := t.IsExist(ctx, filter)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.ErrUserNotFound
	}
	_, err = t.coll.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TourStore) Get(id string) (models.Tour, error) {
	t.logger.Info("repository", "Get", "Tour")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var tour models.Tour
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	ok, err := t.IsExist(ctx, filter)
	if err != nil {
		return models.Tour{}, err
	}
	if !ok {
		return models.Tour{}, errors.ErrUserNotFound
	}
	if err := t.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&tour); err != nil {
		return models.Tour{}, err
	}
	return tour, nil

}

func (t *TourStore) GetAll(page, limit int) ([]models.Tour, int, error) {
	t.logger.Info("repository", "GetAll", "Tour")

	tours := []models.Tour{}
	ctx, cancle := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancle()
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	data, err := t.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}

	for data.Next(ctx) {
		var tour models.Tour
		if err := data.Decode(&tour); err != nil {
			return nil, 0, err
		}
		tours = append(tours, tour)
	}
	count, err := t.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return tours, int(count), nil
}

func (t *TourStore) Update(_id string, updated models.Tour) (models.Tour, error) {
	t.logger.Info("repository", "Update", "Tour")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(_id)
	filter := bson.M{"_id": objId}
	ok, err := t.IsExist(ctx, filter)
	if err != nil {
		return models.Tour{}, err
	}
	if !ok {
		return models.Tour{}, errors.ErrUserNotFound
	} else {
		updated := bson.M{
			"title":        updated.Title,
			"city":         updated.City,
			"address":      updated.Address,
			"photo":        updated.Photo,
			"desc":         updated.Description,
			"price":        updated.Price,
			"maxGroupSize": updated.MaxGroupSize,
			"updatedAt":    time.Now().UTC(),
		}

		_, err := t.coll.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updated})
		if err != nil {
			return models.Tour{}, err
		}
		return t.Get(_id)
	}
}

func (t *TourStore) SearchTour(city string, distance float32, maxgroupsize int) ([]models.Tour, error) {
	t.logger.Info("repository", "SearchTour", "Tour")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	filter := bson.M{"city": city, "distance": bson.M{"$gte": distance}, "maxGroupSize": bson.M{"$gte": maxgroupsize}}
	cur, err := t.coll.Find(ctx, filter)
	if err != nil {
		return []models.Tour{}, err
	}
	tours := []models.Tour{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var tour models.Tour
		if err := cur.Decode(&tour); err != nil {
			continue
		}
		tours = append(tours, tour)
	}
	if err := cur.Err(); err != nil {
		return []models.Tour{}, err
	}
	return tours, nil
}

func (t *TourStore) FeaturedTour() ([]models.Tour, error) {
	t.logger.Info("repository", "FeaturedTour", "Tour")

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	filter := bson.M{"featured": true}
	ok, err := t.IsExist(ctx, filter)
	if err != nil {
		t.logger.Error("repository", "FeaturedTour", err.Error())
		return []models.Tour{}, err
	}
	if !ok {
		t.logger.Error("repository", "FeaturedTour", fmt.Sprintf("%v is not Found", filter))
		return []models.Tour{}, errors.ErrUserNotFound
	}
	cur, err := t.coll.Find(ctx, filter)
	if err != nil {
		t.logger.Error("repository", "FeaturedTour", err.Error())
		return []models.Tour{}, err
	}
	tours := []models.Tour{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		tour := models.Tour{}
		if err := cur.Decode(&tour); err != nil {
			t.logger.Error("repository", "FeaturedTour", fmt.Sprintf("Could not Decode the %v, err : %v", tour, err.Error()))

			continue
		}
		tours = append(tours, tour)
	}

	if cur.Err() != nil {
		t.logger.Error("repository", "FeaturedTour", cur.Err())
		return []models.Tour{}, err
	}
	return tours, nil

}

func (t *TourStore) CountTours() (int, error) {
	t.logger.Info("repository", "CountTours", "Tour")

	opts := options.Count().SetHint("_id_")
	count, err := t.coll.CountDocuments(context.TODO(), bson.D{}, opts)
	if err != nil {
		t.logger.Error("repository", "CountTours", err.Error())

		return 0, err
	}
	return int(count), nil
}
