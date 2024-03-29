package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	TourID     primitive.ObjectID `bson:"tourtId,ref:tourId"`
	Username   string             `bson:"username,required"`
	ReviewText string             `bson:"reviewText,required"`
	Rating     int                `bson:"rating,required"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty"`
}
