package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `json:"name" validate:"required,min=2,max=50"`
	Catagory   string             `json:"catagory" validate:"required"`
	Start_Date time.Time          `json:"start_date" validate:"required"`
	End_Date   time.Time          `json:"end_date" validate:"required"`
	Created_At time.Time          `json:"created_at" validate:"required"`
	Updated_At time.Time          `json:"updated_at" validate:"required"`
	Menu_Id    string            `json:"menu_id" validate:"required"`
}
