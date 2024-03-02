package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,required,unique"`
	Email    string             `bson:"email,required,unique"`
	Password string             `bson:"password,required"`

	Photo     string    `bson:"photo,omitempty"`
	Role      string    `bson:"role,default:user"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
}
