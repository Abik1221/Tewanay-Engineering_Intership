package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/abik1221/Tewanay-Engineering_Intership/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection = database.OpenCollection(database.Client, "order")
var tableCollection = database.OpenCollection(database.Client, "table")

// GetOrders godoc
// @Summary Get all orders
// @Description Retrieve a list of all orders
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {array} bson.M
// @Failure 500 {object} gin.H{"error": string}
// @Router /orders [get]
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get details of an order by order_id
// @Tags order
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 500 {object} gin.H{"error": string}
// @Router /orders/{order_id} [get]
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		order_id := c.Param("order_id")
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Order not found",
			})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with JSON payload; requires valid Table_Id
// @Tags order
// @Accept json
// @Produce json
// @Param order body models.Order true "Order data"
// @Success 200 {object} primitive.InsertOneResult
// @Failure 400 {object} gin.H{"error": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /orders [post]
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		var table models.Table

		// to create an order related to the table we need to check if the table exists
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(order); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		if order.Table_Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Table ID is required",
			})
			return
		}

		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_Id}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Table not found",
			})
			return
		}
		order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_Id = order.ID.Hex()

		sucess, err := orderCollection.InsertOne(ctx, order)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, sucess)
	}
}

// UpdateOrder godoc
// @Summary Update an order by ID
// @Description Update order fields (e.g. order_status) by order_id with JSON payload
// @Tags order
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param order body models.Order true "Order data"
// @Success 200 {object} mongo.UpdateResult
// @Failure 400 {object} gin.H{"error": string}
// @Failure 404 {object} gin.H{"message": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /orders/{order_id} [put]
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		order_Id := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var UpdateObj primitive.D
		if order.Order_Status != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "order_status", Value: order.Order_Status})
		}
		UpdateObj = append(UpdateObj, bson.E{Key: "updated_at", Value: time.Now()})

		upsert := true
		filter := bson.M{"order_id": order_Id}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: UpdateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Order not found",
			})
			return
		}
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "No changes made to the order",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Description Delete an order using order_id
// @Tags order
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} gin.H{"message": string}
// @Failure 404 {object} gin.H{"message": string}
// @Failure 500 {object} gin.H{"error": string}
// @Router /orders/{order_id} [delete]
func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		order_Id := c.Param("order_id")
		filter := bson.M{"order_id": order_Id}
		result, err := orderCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error deleting order",
			})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No order found to be deleted with the given ID. Please change the ID",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
}
