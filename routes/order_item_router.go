package routes

import (
	"github.com/abik1221/Tewanay-Engineering_Intership/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(r *gin.Engine) {
	r.GET("/order_items", controllers.GetOrderItems())
	r.GET("/order_items/:order_item_id", controllers.GetOrderItem())
	r.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrderId())
	r.POST("/order_items", controllers.CreateOrderItem())
	r.PATCH("/order_items/:order_item_id", controllers.UpdateOrderItem())
	r.DELETE("/order_items/:order_item_id", controllers.DeleteOrderItem())
}
