package main

import (
	"os"

	"github.com/abik1221/Tewanay-Engineering_Intership/middlewares"
	"github.com/abik1221/Tewanay-Engineering_Intership/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/abik1221/Tewanay-Engineering_Intership/database"
)

var foods *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	routes.UserRoutes(router)
	router.Use(middlewares.AuthMiddleware())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)
	routes.OrderItemRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
