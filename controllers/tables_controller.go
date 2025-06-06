package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTables godoc
// @Summary      Get all tables
// @Description  Retrieve a list of all tables
// @Tags         Tables
// @Produce      json
// @Success      200  {array}   models.Table
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /tables [get]
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var tables []models.Table
		cursor, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Food not found"})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var table models.Table
			if err := cursor.Decode(&table); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding table"})
				return
			}
			tables = append(tables, table)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
			return
		}
		c.JSON(http.StatusOK, tables)

	}
}

// GetTable godoc
// @Summary      Get a table by ID
// @Description  Retrieve a single table by its table_id
// @Tags         Tables
// @Produce      json
// @Param        table_id   path      string  true  "Table ID"
// @Success      200  {object}  models.Table
// @Failure      404  {object}  gin.H{"error": string}
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /tables/{table_id} [get]
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableId := c.Param("table_id")
		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table not found"})
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

// CreateTable godoc
// @Summary      Create a new table
// @Description  Create a new table with JSON input
// @Tags         Tables
// @Accept       json
// @Produce      json
// @Param        table  body      models.Table  true  "Table Data"
// @Success      200  {object}  primitive.InsertOneResult
// @Failure      400  {object}  gin.H{"error": string}
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /tables [post]
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		if err := c.ShouldBindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		table.ID = primitive.NewObjectID()
		table.Table_Id = table.ID.Hex()
		table.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating table"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// UpdateTable godoc
// @Summary      Update an existing table
// @Description  Update a table by table_id with JSON input
// @Tags         Tables
// @Accept       json
// @Produce      json
// @Param        table_id  path      string       true  "Table ID"
// @Param        table     body      models.Table true  "Updated Table Data"
// @Success      200  {object}  primitive.UpdateResult
// @Failure      400  {object}  gin.H{"error": string}
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /tables/{table_id} [put]
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		tableId := c.Param("table_id")
		if err := c.ShouldBindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		filter := bson.M{"table_id": tableId}
		update := bson.M{"$set": table}
		result, err := tableCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating table"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// DeleteTable godoc
// @Summary      Delete a table
// @Description  Delete a table by table_id
// @Tags         Tables
// @Produce      json
// @Param        table_id  path      string  true  "Table ID"
// @Success      200  {object}  primitive.DeleteResult
// @Failure      404  {object}  gin.H{"error": string}
// @Failure      500  {object}  gin.H{"error": string}
// @Router       /tables/{table_id} [delete]
func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableId := c.Param("table_id")
		filter := bson.M{"table_id": tableId}
		result, err := tableCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting table"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
