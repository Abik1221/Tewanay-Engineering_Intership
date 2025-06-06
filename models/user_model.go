package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	First_Name string             `bson:"first_name" json:"first_name" validate:"required"`
	Last_Name  string             `bson:"last_name" json:"last_name" validate:"required"`
	Password  string             `json:"password" validate:"required"`
	Email     string             `json:"email" validate:"required,email"`
	Phone     string             `json:"phone" validate:"required"`
	Role      string             `json:"role" validate:"required,oneof=admin user"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
