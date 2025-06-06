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
// @Summary Get all tables
// @Description Retrieve all tables from the database
// @Tags tables
// @Accept json
// @Produce json
// @Success 200 {array} models.Table
// @Failure 500 {object} object "Internal Server Error"
// @Router /tables [get]
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
// @Summary Get a single table
// @Description Retrieve a table by its ID
// @Tags tables
// @Accept json
// @Produce json
// @Param table_id path string true "Table ID"
// @Success 200 {object} models.Table
// @Failure 404 {object} object "Table not found"
// @Failure 500 {object} object "Internal Server Error"
// @Router /tables/{table_id} [get]
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableId := c.Param("table_id")
		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table not found"})
		}
		c.JSON(http.StatusOK, table)
	}
}

// CreateTable godoc
// @Summary Create a new table
// @Description Add a new table to the database
// @Tags tables
// @Accept json
// @Produce json
// @Param table body models.Table true "Table data"
// @Success 200 {object} object "Insertion result"
// @Failure 400 {object} object "Invalid input"
// @Failure 500 {object} object "Error creating table"
// @Router /tables [post]
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
// @Summary Update a table
// @Description Update an existing table's information
// @Tags tables
// @Accept json
// @Produce json
// @Param table_id path string true "Table ID"
// @Param table body models.Table true "Updated table data"
// @Success 200 {object} object "Update result"
// @Failure 400 {object} object "Invalid input"
// @Failure 500 {object} object "Error updating table"
// @Router /tables/{table_id} [put]
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
// @Summary Delete a table
// @Description Remove a table from the database
// @Tags tables
// @Accept json
// @Produce json
// @Param table_id path string true "Table ID"
// @Success 200 {object} object "Delete result"
// @Failure 404 {object} object "Table not found"
// @Failure 500 {object} object "Error deleting table"
// @Router /tables/{table_id} [delete]
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
