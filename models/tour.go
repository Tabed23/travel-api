package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tour struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title,required,unique"`
	City         string             `bson:"city,required"`
	Address      string             `bson:"address,required"`
	Distance     int                `bson:"distance,required"`
	Photo        string             `bson:"photo,required"`
	Description  string             `bson:"desc,required"`
	Price        float64            `bson:"price,required"`
	MaxGroupSize int                `bson:"maxGroupSize,required"`

	Reviews   []primitive.ObjectID `bson:"reviews,ref:Reviews"`
	Featured  bool                 `bson:"featured,default:false"`
	CreatedAt time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty"`
}
