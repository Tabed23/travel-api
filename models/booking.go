package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserEmail string             `bson:"user_email"`
	TourName  string             `bson:"tour_name"`
	GuestSize int64              `bson:"guest_size"`
	Phone     int64              `bson:"phone"`
	CreatedAt time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty"`
}


type UpdateBooking struct {
	TourName  string             `bson:"tour_name"`
	GuestSize int64              `bson:"guest_size"`
	Phone     int64              `bson:"phone"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty"`
}