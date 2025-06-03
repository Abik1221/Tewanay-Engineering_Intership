package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ordered_Item struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Menu_Id    string             `json:"menu_id" validate:"required"`
	Food_Id    string             `json:"food_id" validate:"required"`
	Order_Id   string             `json:"order_id" validate:"required"`
	Quantity   int                `json:"quantity" validate:"required"`
	Price      float64            `json:"price" validate:"required"`
	Created_At time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_At time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
