package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Order_Id     string             `json:"order_id" validate:"required"`
	Table_Id     string             `json:"table_id" validate:"required"`
	Order_Status string             `json:"order_status" validate:"required"`
	Created_At   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_At   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
