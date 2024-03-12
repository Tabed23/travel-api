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

type BookingStore struct {
	b      mongo.Collection
	u      mongo.Collection
	logger *slog.Logger
}

func NewBookStore(c mongo.Collection, l *slog.Logger) *BookingStore {
	return &BookingStore{b: c, logger: l}
}

// GetBooking By ID
func (b *BookingStore) GetBooking(id string) (models.Booking, error) {
	b.logger.Info("repository", "Get", "Booking")
	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	ok, err := b.IsExist(ctx, filter)
	if err != nil {
		b.logger.Error("repository", "IsExist", err.Error())
		return models.Booking{}, err
	}
	if !ok {
		b.logger.Error("repository", "Get", "Booking not found")
		return models.Booking{}, errors.ErrUserNotFound
	}

	var booking models.Booking
	if err := b.b.FindOne(ctx, objId).Decode(&booking); err != nil {
		b.logger.Error("repository", "Get", err.Error())
		return models.Booking{}, err
	}
	b.logger.Info("repository", "Get", fmt.Sprintf("Booking found: %v", booking))

	return booking, nil
}

// CreateBooking Create a new booking
func (b *BookingStore) CreateBooking(userId string, book models.Booking) (models.Booking, error) {
	b.logger.Info("repository", "Create", "Booking")

	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		b.logger.Error("repository", "CreateReviw", err.Error())
		return models.Booking{}, err
	}
	usr := models.User{}
	if err := b.u.FindOne(ctx, bson.M{"_id": objId}).Decode(&usr); err != nil {
		b.logger.Error("repository", "Create", err.Error())
		return models.Booking{}, err
	}
	book.ID = primitive.NewObjectID()
	book.UserID = usr.ID
	book.UserEmail = usr.Email
	book.CreatedAt = time.Now().UTC()
	book.UpdatedAt = time.Now().UTC()
	if _, err := b.b.InsertOne(ctx, &book); err != nil {
		b.logger.Error("repository", "Create", err.Error())

		return models.Booking{}, err
	}
	b.logger.Info("repository", "Create", fmt.Sprintf("Booking created: %v", book))
	return book, nil
}

// Delete Booking Document
func (b *BookingStore) DeleteBook(id string) (bool, error) {
	b.logger.Info("repository", "Delete", "Booking")

	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	ok, err := b.IsExist(ctx, filter)
	if err != nil {
		b.logger.Error("repository", "IsExist", err.Error())
		return false, err
	}
	if !ok {
		b.logger.Warn("repository", "IsExist", "not found")

		return false, errors.ErrUserNotFound
	}

	if _, err := b.b.DeleteOne(ctx, bson.M{"_id": objId}); err != nil {
		b.logger.Error("repository", "Delete", err.Error())
		return false, err
	}
	b.logger.Info("repository", "Delete", fmt.Sprintf("Booking deleted: %v", objId))
	return true, nil
}

func (b *BookingStore) Update(id string, bookUpdate models.UpdateBooking) (models.Booking, error) {
	b.logger.Info("repository", "Update", "Booking")

	ctx, cancle := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancle()
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}
	ok, err := b.IsExist(ctx, filter)
	if err != nil {
		b.logger.Error("repository", "Update", err.Error())

		return models.Booking{}, err
	}
	if !ok {
		b.logger.Warn("repository", "Update", fmt.Sprintf("Booking not found: %v", id))

		return models.Booking{}, errors.ErrUserNotFound
	}

	update := bson.M{
		"tour_name":  bookUpdate.TourName,
		"guest_size": bookUpdate.GuestSize,
		"phone":      bookUpdate.Phone,
		"updatedAt":  time.Now().UTC(),
	}
	_, err = b.b.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		b.logger.Error("repository", "Update", err.Error())

		return models.Booking{}, err
	}
	var booking models.Booking
	if err = b.b.FindOne(ctx, bson.M{"_id": objId}).Decode(&booking); err != nil {
		b.logger.Error("repository", "Update", err.Error())

		return models.Booking{}, err
	}

	b.logger.Info("repository", "Update", fmt.Sprintf("Booking update: %v", booking))

	return booking, nil
}

// Get All Bookings Documemt0s
func (b *BookingStore) GetAll(page, limit int) ([]models.Booking, int, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	skip := (page - 1) * limit

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	bookings := []models.Booking{}

	cur, err := b.b.Find(ctx, bson.M{}, opts)
	if err != nil {
		b.logger.Error("repository", "GetAll", err.Error())
		return []models.Booking{}, 0, err
	}
	for cur.Next(ctx) {
		var booking models.Booking
		if err := cur.Decode(&booking); err != nil {
			b.logger.Error("repository", "GetAll", err.Error())
			return []models.Booking{}, 0, err
		}
		bookings = append(bookings, booking)
	}
	if err := cur.Err(); err != nil {
		b.logger.Error("repository", "GetAll", err.Error())
		return []models.Booking{}, 0, err
	}
	count, err := b.b.CountDocuments(ctx, bson.M{})
	if err != nil {
		b.logger.Error("repository", "GetAll", err.Error())
		return []models.Booking{}, 0, err
	}
	return bookings, int(count), err
}

// Count Documents
func (b *BookingStore) Count() (int, error) {
	opts := options.Count().SetHint("_id_")
	count, err := b.b.CountDocuments(context.Background(), opts)
	if err != nil {
		b.logger.Error("respsitory", "count", err.Error())
		return 0, err
	}
	b.logger.Info("respsitory", "count", fmt.Sprintf("%v", count))
	return int(count), nil
}

// IsExit Documents
func (b *BookingStore) IsExist(ctx context.Context, filter primitive.M) (bool, error) {
	count, err := b.b.CountDocuments(ctx, filter)
	if err != nil {
		b.logger.Error("repository", "ISExist", err.Error())
		return false, err
	}

	if count > 0 {
		b.logger.Info("repository", "ISExist", "Found")
		return true, nil
	}
	b.logger.Warn("repository", "ISExits", "Not Found")
	return false, nil

}
