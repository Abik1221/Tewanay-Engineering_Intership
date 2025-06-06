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


// @Summary      List all order items
// @Description  Retrieve all order items in the system
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Success      200  {array}   primitive.M
// @Failure      404  {object}  object  "No order items found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /order-items [get]
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

// @Summary      Get an order item by ID
// @Description  Fetch a single order item by its ID
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_item_id  path  string  true  "Order Item ID"
// @Success      200  {object}  primitive.M
// @Failure      404  {object}  object  "Order item not found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /order-items/{order_item_id} [get]
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

// @Summary      Get order items by order ID
// @Description  Retrieve all order items for a specific order
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_id  path  string  true  "Order ID"
// @Success      200  {array}   primitive.M
// @Failure      404  {object}  object  "No order items found for this order"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /order-items/by-order/{order_id} [get]
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

// @Summary      Create a new order item
// @Description  Add a new order item to the database
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        request  body  models.Ordered_Item  true  "Order item data"
// @Success      200  {object}  object  "MongoDB insert result"
// @Failure      400  {object}  object  "Invalid input"
// @Failure      500  {object}  object  "Error creating order item"
// @Router       /order-items [post]
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

// @Summary      Update an order item
// @Description  Modify an existing order item
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_item_id  path  string                true  "Order Item ID"
// @Param        request        body  models.Ordered_Item  true  "Order item data"
// @Success      200  {object}  object  "MongoDB update result"
// @Failure      400  {object}  object  "Invalid input"
// @Failure      404  {object}  object  "Order item not found"
// @Failure      500  {object}  object  "Error updating order item"
// @Router       /order-items/{order_item_id} [put]
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


// @Summary      Delete an order item
// @Description  Remove an order item by ID
// @Tags         order-items
// @Accept       json
// @Produce      json
// @Param        order_item_id  path  string  true  "Order Item ID"
// @Success      200  {object}  object  "message: Order item deleted successfully"
// @Failure      404  {object}  object  "Order item not found"
// @Failure      500  {object}  object  "Error deleting order item"
// @Router       /order-items/{order_item_id} [delete]
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
