package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Table_Id   string             `json:"table_id" validate:"required"`
	Table_Name string             `json:"table_name" validate:"required"`
	Created_At time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_At time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
