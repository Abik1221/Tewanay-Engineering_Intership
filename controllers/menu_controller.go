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

var menuCollection = database.OpenCollection(database.Client, "menu")

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menu_id := c.Param("menu_id")
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menu_id}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Menu not found",
			})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch menus",
			})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(menu); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}
		menu.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_Id = menu.ID.Hex()

		sucess, err := menuCollection.InsertOne(ctx, menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, sucess)

	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		menu_id := c.Param("menu_id")
		filter := bson.M{"menu_id": menu_id}

		var UpdateObj primitive.D

		if !menu.Start_Date.IsZero() && !menu.End_Date.IsZero() {

			if !inTimeSpan(menu.Start_Date, menu.End_Date, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Start date must be before end date and both must be in the future",
				})
				return
			}

			UpdateObj = append(UpdateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			UpdateObj = append(UpdateObj, bson.E{Key: "end_date", Value: menu.End_Date})

		}

		if menu.Name != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "name", Value: menu.Name})
		}

		if menu.Catagory != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "catagory", Value: menu.Catagory})
		}

		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		UpdateObj = append(UpdateObj, bson.E{Key: "updated_at", Value: menu.Updated_At})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", UpdateObj},
			},
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

func inTimeSpan(start, end, now time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menu_id := c.Param("menu_id")
		filter := bson.M{"menu_id": menu_id}
		result, err := menuCollection.DeleteOne(ctx, filter)
		if err != nil {
			log.Println("Error deleting menu:", err)
			return
		}
		if result.DeletedCount == 0 {
			log.Println("No menu found with the given ID")
			return
		}
		log.Println("Menu deleted successfully")
		c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
	}
}
