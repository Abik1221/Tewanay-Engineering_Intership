package routes

import (
	"github.com/abik1221/Tewanay-Engineering_Intership/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controllers.GetUsers())
	r.GET("users/:user_id", controllers.GetUser())
	r.POST("/users/signup", controllers.Signup())
	r.POST("/users/login", controllers.Login())
}
