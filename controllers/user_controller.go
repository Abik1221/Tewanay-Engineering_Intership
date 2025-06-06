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
          var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		  defer cancel()

		  userId := c.Param("user_id")

		  var user models.User

		  err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		  if err != nil {
			c.Json(http.StatusInternalServerError, gin.H{"error": err.Error()})
		  }

		  c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
         var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
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
		 user.Password = &password

		 count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		 if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		 }
		 if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
			return
		 }

		 user.CreatedAt = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		 user.UpdatedAt = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		 user.ID = primitive.NewObjectID()
		 user.user_id = user.ID.Hex()

		token, refrest_tokens, _ :=  helper.GenerateAllTokens(* user.Email, *user.First_Name, *user.Last_Name, user.user_id)
		 user.Token = &token
		 user.Refresh_Token = &refrest_tokens

        result, err :=userCollection.InsertOne(ctx, user)
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
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
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		PasswordIsValid := VerifyPassword(user.Password, *user.Password)
		if !PasswordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refresh_token, _ := helper.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, user.user_id)
		helper.UpdateAllTokens(token, refresh_token, user.User_id)

		user.Token = &token
		user.Refresh_Token = &refresh_token

		c.JSON(http.StatusOK, user)
	}
}

func HashPassward(User_Password string) string {
	bytes, err := bycrypt.GenerateFromPassword([]byte(User_Password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassward(Intered_password, hash string) (bool, string) {
	err := bycrypt.CompareHashAndPassword([]byte(hash), []byte(Intered_password))
	if err != nil {
		return false, "Invalid password"
	}
	return true, ""
}
