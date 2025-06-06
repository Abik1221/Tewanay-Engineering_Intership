package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/abik1221/Tewanay-Engineering_Intership/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var invoiceCollection = database.OpenCollection(database.Client, "invoices")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoices []models.Invoice
		cursor, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invoices not found"})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var invoice models.Invoice
			if err := cursor.Decode(&invoice); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding invoice"})
				return
			}
			invoices = append(invoices, invoice)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
			return
		}
		c.JSON(http.StatusOK, invoices)

	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invoice not found"})
			return
		}
		c.JSON(http.StatusOK, invoice)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(invoice); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		invoice.Invoice_Id = primitive.NewObjectID().Hex()
		invoice.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		sucess, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating invoice"})
			return
		}
		c.JSON(http.StatusOK, sucess)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"invoice_id": invoiceId}
		invoice.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		update := bson.M{
			"$set": invoice,
		}
		result, err := invoiceCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating invoice"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		invoiceId := c.Param("invoice_id")
		filter := bson.M{"invoice_id": invoiceId}
		result, err := invoiceCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting invoice"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
