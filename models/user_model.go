package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type user struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}
