package routes

import (
	"github.com/abik1221/Tewanay-Engineering_Intership/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	r.GET("/orders", controllers.GetOrders())
	r.GET("/orders/:order_id", controllers.GetOrder())
	r.POST("/orders", controllers.CreateOrders())
	r.PATCH("/orders/:order_id", controllers.UpdateOrder())
	r.DELETE("orders/:order_id", controllers.DeleteOrder())
}
