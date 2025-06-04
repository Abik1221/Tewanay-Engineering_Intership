package controllers

import (
	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/gin-gonic/gin"
)

var userCollection = database.OpenCollection((dataase.Client, "user"))

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
       ctx, cancel := context.WithTimeout(Context.Background(), 100 * time.Second)
	   defer cancel()
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassward(User_Password string) string {
	return ""
}

func VerifyPassward(Intered_password, hash string) bool {
	return false
}
