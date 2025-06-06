package helpers

import (
	"context"
	"os"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	User_id    string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var secretKey = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, user_id string) (signedToken string, refresh_token string, err error) {

	claims := &signedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		User_id:    user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}
	refresh_claims := &signedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		User_id:    user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(72))),
		},
	}

	tokens, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refresh_tokens, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return tokens, refresh_tokens, nil

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, user_id string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updatedat := time.Now().Format(time.RFC3339)

	updateObj = append(updateObj, bson.E{"updated_at", updatedat})

	upsert := true

	filter := bson.M{"user_id": user_id}
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opts)
	if err != nil {
		return
	}

}

func ValidateAllTokens(signedToken string) (claims *signedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*signedDetails)
	if !ok || !token.Valid {
		msg = "the token is invalid"
		return nil, msg
	}

	if claims.ExpiresAt.Time.Before(time.Now().Local()) {
		msg = "the token is expired"
		return nil, msg
	}

	return claims, msg

}
