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

// GetInvoices godoc
// @Summary      Get all invoices
// @Description  Retrieves a list of all invoices from the database
// @Tags         invoices
// @Produce      json
// @Success      200  {array}   models.Invoice
// @Failure      500  {object}  gin.H{"error": "Invoices not found"}
// @Router       /invoices [get]
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

// GetInvoice godoc
// @Summary      Get an invoice by ID
// @Description  Retrieves a single invoice by its invoice_id
// @Tags         invoices
// @Produce      json
// @Param        invoice_id  path      string  true  "Invoice ID"
// @Success      200         {object}  models.Invoice
// @Failure      500         {object}  gin.H{"error": "Invoice not found"}
// @Router       /invoices/{invoice_id} [get]
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

// CreateInvoice godoc
// @Summary      Create a new invoice
// @Description  Creates a new invoice document in the database
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        invoice  body      models.Invoice  true  "Invoice data"
// @Success      200      {object}  primitive.InsertOneResult
// @Failure      400      {object}  gin.H{"error": "Invalid input"}
// @Failure      500      {object}  gin.H{"error": "Error creating invoice"}
// @Router       /invoices [post]
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
		success, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating invoice"})
			return
		}
		c.JSON(http.StatusOK, success)
	}
}

// UpdateInvoice godoc
// @Summary      Update an invoice by ID
// @Description  Updates an existing invoice's data by invoice_id
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Param        invoice_id  path      string         true  "Invoice ID"
// @Param        invoice     body      models.Invoice  true  "Updated invoice data"
// @Success      200        {object}  mongo.UpdateResult
// @Failure      400        {object}  gin.H{"error": "Invalid input"}
// @Failure      500        {object}  gin.H{"error": "Error updating invoice"}
// @Router       /invoices/{invoice_id} [put]
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

// DeleteInvoice godoc
// @Summary      Delete an invoice by ID
// @Description  Deletes an invoice document from the database by invoice_id
// @Tags         invoices
// @Produce      json
// @Param        invoice_id  path      string  true  "Invoice ID"
// @Success      200         {object}  mongo.DeleteResult
// @Failure      500         {object}  gin.H{"error": "Error deleting invoice"}
// @Router       /invoices/{invoice_id} [delete]
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
