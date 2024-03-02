package store

import (
	"context"
	"time"

	"github.com/tabed23/travel-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TourStore struct {
	coll mongo.Collection
}

func NewTourStore(c mongo.Collection) *TourStore {
	return &TourStore{coll: c}
}
func (t *TourStore) CreateTour(tour models.Tour) (models.Tour, error) {
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

func (t *TourStore) IsExist(ctx context.Context, value, data string) (bool, error) {
	count, err := t.coll.CountDocuments(ctx, bson.M{value: data})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil

}

func (t *TourStore) DeleteTour(_id string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(_id)
	_, err := t.coll.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TourStore) Get(id string) (models.Tour, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var tour models.Tour
	objId, _ := primitive.ObjectIDFromHex(id)
	if err := t.coll.FindOne(ctx, bson.M{"_id": objId}).Decode(&tour); err != nil {
		return models.Tour{}, err
	}
	return tour, nil

}

func (t *TourStore) GetAll(page, limit int) ([]models.Tour, int, error) {
	tours := []models.Tour{}
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	ok, err := t.IsExist(ctx, _id, "_id")
	if err != nil {
		return models.Tour{}, err
	}
	if ok {
		objId, _ := primitive.ObjectIDFromHex(_id)
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
	return models.Tour{}, nil
}

func (t *TourStore) SearchTour(city string, distance float32, maxgroupsize int) ([]models.Tour, error) {
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

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	filter := bson.M{"featured": true}
	cur, err := t.coll.Find(ctx, filter)
	if err != nil {
		return []models.Tour{}, err
	}
	tours := []models.Tour{}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		tour := models.Tour{}
		if err := cur.Decode(&tour); err != nil {
			continue
		}
		tours = append(tours, tour)
	}

	if cur.Err() != nil {
		return []models.Tour{}, err
	}
	return tours, nil

}

func (t *TourStore) CountTours() (int, error) {
	opts := options.Count().SetHint("_id_")
	count, err := t.coll.CountDocuments(context.TODO(), bson.D{}, opts)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
