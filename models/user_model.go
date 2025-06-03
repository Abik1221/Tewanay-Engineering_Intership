package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `json:"username" validate:"required"`
	Password  string             `json:"password" validate:"required"`
	Email     string             `json:"email" validate:"required,email"`
	Role      string             `json:"role" validate:"required,oneof=admin user"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
