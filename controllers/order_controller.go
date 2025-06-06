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

// @Summary      List all orders
// @Description  Retrieve a list of all orders in the system
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   bson.M
// @Failure      500  {object}  object  "Internal server error"
// @Router       /orders [get]
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

// @Summary      Get an order by ID
// @Description  Fetch a single order by its unique ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order_id  path  string  true  "Order ID"
// @Success      200  {object}  models.Order
// @Failure      404  {object}  object  "Order not found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /orders/{order_id} [get]
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		order_id := c.Param("order_id")
		var order models.Order
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": order_id}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Order Not found",
			})
		}
		c.JSON(http.StatusOK, order)
	}
}

// @Summary      Create a new order
// @Description  Add a new order to the database (requires valid table_id)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request  body  models.Order  true  "Order data (must include table_id)"
// @Success      200  {object}  object  "MongoDB insert result"
// @Failure      400  {object}  object  "Invalid input or missing table_id"
// @Failure      404  {object}  object  "Table not found"
// @Failure      500  {object}  object  "Error creating order"
// @Router       /orders [post]
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

// @Summary      Create a new order
// @Description  Add a new order to the database (requires valid table_id)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request  body  models.Order  true  "Order data (must include table_id)"
// @Success      200  {object}  object  "MongoDB insert result"
// @Failure      400  {object}  object  "Invalid input or missing table_id"
// @Failure      404  {object}  object  "Table not found"
// @Failure      500  {object}  object  "Error creating order"
// @Router       /orders [post]
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

// @Summary      Delete an order
// @Description  Remove an order by ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order_id  path  string  true  "Order ID to delete"
// @Success      200  {object}  object  "MongoDB delete result"
// @Failure      404  {object}  object  "Order not found"
// @Failure      500  {object}  object  "Error deleting order"
// @Router       /orders/{order_id} [delete]
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
				"message": "no order found to be deleted with the given id please cange the ID",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
