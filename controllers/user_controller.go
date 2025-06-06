package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/abik1221/Tewanay-Engineering_Intership/helpers"
	"github.com/abik1221/Tewanay-Engineering_Intership/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var users []models.User

		results, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer results.Close(ctx)

		for results.Next(ctx) {
			var user models.User
			if err := results.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validation_err := validate.Struct(user)
		if validation_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validation_err.Error(),
			})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		password := HashPassward(user.Password)
		user.Password = password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refrest_tokens, _ := helpers.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refrest_tokens

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var Found_user models.User
		var user models.User
		if err := c.BindJSON(&Found_user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validation_err := validate.Struct(Found_user)
		if validation_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validation_err.Error(),
			})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": Found_user.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		PasswordIsValid, _ := VerifyPassward(Found_user.Password, user.Password)
		if !PasswordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refresh_token, _ := helpers.GenerateAllTokens(Found_user.Email, Found_user.First_Name, Found_user.Last_Name, Found_user.User_id)
		helpers.UpdateAllTokens(token, refresh_token, Found_user.User_id)

		Found_user.Token = &token
		Found_user.Refresh_Token = &refresh_token

		c.JSON(http.StatusOK, Found_user)
	}
}
func HashPassward(User_Password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(User_Password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassward(Intered_password, hash string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(Intered_password))
	if err != nil {
		return false, "Invalid password"
	}
	return true, ""
}
