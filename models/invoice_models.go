package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Invoice_Id   string   `json:"invoice_id" validate:"required"`
	Order_Id   string  `json:"order_id" validate:"required"`
	Payment_Method  *string `json:"payment_method" validate:"required"`
	Payment_Status  *string `json:"payment_status" validate:"required"`
	Payment_Due_Date  time.Time `json:"payment_due_date" validate:"required"`
	Created_At  time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_At time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}