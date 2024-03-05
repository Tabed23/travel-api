package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name,required"`
	Lastname  string             `bson:"last_name,required"`
	UserName  string             `bson:"user_name,required"`
	Email     string             `bson:"email,required,unique"`
	Password  string             `bson:"password,required"`

	Photo     string    `bson:"photo,omitempty"`
	Role      string    `bson:"role,default:user"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
}

type UserLogin struct {
	Email    string `bson:"email,required,unique"`
	Password string `bson:"password,required"`
}

type UserUpdate struct {
	FirstName string `bson:"first_name,required"`
	Lastname  string `bson:"last_name,required"`
	Email     string `bson:"email,required,unique"`
	Password  string `bson:"password,required"`

	Photo string `bson:"photo,omitempty"`
	Role  string `bson:"role,default:user"`
}
