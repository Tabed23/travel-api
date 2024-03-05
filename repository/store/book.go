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

type BookingStore struct {
	b mongo.Collection
}

func NewBookStore(c mongo.Collection) *BookingStore {
	return &BookingStore{b: c}
}

func (s *BookingStore) GetBooking(id string) (models.Booking, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	var booking models.Booking
	if err := s.b.FindOne(ctx, objId).Decode(&booking); err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (s *BookingStore) CreateBooking(book models.Booking) (models.Booking, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	book.ID = primitive.NewObjectID()
	book.CreatedAt = time.Now().UTC()
	book.UpdatedAt = time.Now().UTC()
	if _, err := s.b.InsertOne(ctx, &book); err != nil {
		return models.Booking{}, err
	}

	return book, nil
}

func (s *BookingStore) DeleteBook(id string) (bool, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	if _, err := s.b.DeleteOne(ctx, bson.M{"_id": objId}); err != nil {
		return false, err
	}
	return true, nil
}

func (s *BookingStore) Update(id string, bookUpdate models.UpdateBooking) (models.Booking, error) {
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	update := bson.M{
		"tour_name":  bookUpdate.TourName,
		"guest_size": bookUpdate.GuestSize,
		"phone":      bookUpdate.Phone,
		"updatedAt":  time.Now().UTC(),
	}
	_, err := s.b.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return models.Booking{}, err
	}
	var booking models.Booking
	if err = s.b.FindOne(ctx, bson.M{"_id": objId}).Decode(&booking); err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (s *BookingStore) GetAll(page, limit int) ([]models.Booking, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	bookings := []models.Booking{}

	cur, err := s.b.Find(ctx, bson.M{}, opts)
	if err != nil {
		return []models.Booking{}, 0, err
	}
	for cur.Next(ctx) {
		var booking models.Booking
		if err := cur.Decode(&booking); err != nil {
			return []models.Booking{}, 0, err
		}
		bookings = append(bookings, booking)
	}
	if err := cur.Err(); err != nil {
		return []models.Booking{}, 0, err
	}
	count, err := s.b.CountDocuments(ctx, bson.M{})
	if err != nil {
		return []models.Booking{}, 0, err
	}
	return bookings, int(count), err
}
