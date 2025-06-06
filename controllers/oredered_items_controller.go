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
	"go.mongodb.org/mongo-driver/mongo"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order_items")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {

		var orederItems []primitive.M
		orederItems, err := itemsByOrder("")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(orederItems) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No order items found"})
			return
		}
		c.JSON(http.StatusOK, orederItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderItemId := c.Param("order_item_id")
		var orderItem primitive.M
		err := orderItemCollection.FindOne(ctx, bson.M{"_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order item not found"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func GetOrderItemsByOrderId() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		var orederItems []primitive.M
		orederItems, err := itemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(orederItems) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No order items found for this order"})
			return
		}
		c.JSON(http.StatusOK, orederItems)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.Ordered_Item
		if err := c.ShouldBindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		orderItem.ID = primitive.NewObjectID()
		orderItem.Order_Id = orderItem.ID.Hex()
		orderItem.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, err := orderItemCollection.InsertOne(ctx, orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order item"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.Ordered_Item
		orderItemId := c.Param("order_item_id")
		if err := c.ShouldBindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		orderItem.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"order_item_id": orderItemId}
		update := bson.M{"$set": orderItem}
		result, err := orderItemCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order item"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Order item not found"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}
		result, err := orderItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order item"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
	}
}

func itemsByOrder(id string) (OrederItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "null"}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
	projectStage := bson.D{{"$project", bson.D{{"_id", 0}, {"order_items", 1}}}}
	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		return nil, err
	}
	var orederItems []primitive.M
	if err = result.All(ctx, &orederItems); err != nil {
		return nil, err
	}
	return orederItems, nil
}
