package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Food_Name        string             `json:"food_name" validate:"required,min=2,max=50"`
	Food_Price       *float64           `json:"food_price" validate:"required"`
	Food_Description string             `json:"food_description" validate:"required"`
	Food_Image       string             `json:"food_image" validate:"required"`
	Created_AT       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_AT       time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Food_Id          *string            `json:"food_id" validate:"required"`
	Menu_Id          *string            `json:"menu_id" validate:"required"`
}