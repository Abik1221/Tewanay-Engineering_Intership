package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/abik1221/Tewanay-Engineering_Intership/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()
var foodCollection = database.OpenCollection(database.Client, "food")

// @Summary      Get a single food item
// @Description  Fetch food details by its unique ID
// @Tags         foods
// @Accept       json
// @Produce      json
// @Param        food_id  path  string  true  "Food ID"
// @Success      200  {object}  models.Food
// @Failure      404  {object}  object  "Food not found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /foods/{food_id} [get]
func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodID := c.Param("food_id")

		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Food not found"})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

// @Summary      List all foods (paginated)
// @Description  Retrieve a paginated list of food items with optional query params
// @Tags         foods
// @Accept       json
// @Produce      json
// @Param        page          query  int     false  "Page number (default: 1)"
// @Param        recordPerPage query  int     false  "Items per page (default: 10)"
// @Param        startIndex    query  int     false  "Custom start index (overrides page)"
// @Success      200  {object}  []bson.M
// @Failure      500  {object}  object  "Internal server error"
// @Router       /foods [get]
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <= 0 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "null"},
				{Key: "totalCount", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "data", Value: bson.D{
					{Key: "$push", Value: "$$ROOT"},
				}},
			}},
		}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "totalCount", Value: 1},
					{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				},
			},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, allFoods)
	}
}

// @Summary      Create a new food item
// @Description  Add a new food entry to the database
// @Tags         foods
// @Accept       json
// @Produce      json
// @Param        request  body  models.Menu  true  "Food data (Note: Uses Menu model for binding)"
// @Success      200  {object}  object  "MongoDB insert result"
// @Failure      400  {object}  object  "Invalid input or validation error"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /foods [post]
func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input data",
			})
			return
		}

		if err := validate.Struct(food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		food.Created_AT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_AT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		foodIdHex := food.ID.Hex()
		food.Food_Id = &foodIdHex
		var num = toFixed(*food.Food_Price, 2)
		food.Food_Price = &num

		menu.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := foodCollection.InsertOne(ctx, food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating a food"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary      Update a food item
// @Description  Modify food details by ID (partial updates supported)
// @Tags         foods
// @Accept       json
// @Produce      json
// @Param        food_id  path  string       true  "Food ID to update"
// @Param        request  body  models.Food  true  "Fields to update (all optional)"
// @Success      200  {object}  object  "MongoDB update result"
// @Failure      400  {object}  object  "Invalid input"
// @Failure      404  {object}  object  "Menu ID not found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /foods/{food_id} [patch]
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		var food models.Food

		food_Id := c.Param("food_id")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var UpdateObj primitive.D
		if food.Food_Name != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "food_name", Value: food.Food_Name})
		}

		if food.Food_Price != nil {
			UpdateObj = append(UpdateObj, bson.E{Key: "food_price", Value: food.Food_Price})
		}

		if food.Food_Image != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "food_image", Value: food.Food_Image})
		}

		if food.Menu_Id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_Id}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				},
				)
				return
			}
			UpdateObj = append(UpdateObj, bson.E{Key: "menu_id", Value: food.Menu_Id})
		}

		if food.Food_Description != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "food_description", Value: food.Food_Description})
		}

		food.Updated_AT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		UpdateObj = append(UpdateObj, bson.E{Key: "updated_at", Value: food.Updated_AT})

		upsert := true
		filter := bson.M{"food_id": food_Id}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: UpdateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

// @Summary      Delete a food item
// @Description  Remove a food entry by ID
// @Tags         foods
// @Accept       json
// @Produce      json
// @Param        food_id  path  string  true  "Food ID to delete"
// @Success      200  {object}  object  "message: Food deleted successfully"
// @Failure      404  {object}  object  "Food not found"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /foods/{food_id} [delete]
func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		food_Id := c.Param("food_id")

		filter := bson.M{"food_id": food_Id}

		result, err := foodCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error deleting food",
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Food not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Food deleted successfully",
		})
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
