package helpers

import (
	"os"
	"time"

	"github.com/abik1221/Tewanay-Engineering_Intership/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type signedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	User_id    string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var secretKey = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, user_id string) (signedToken string, refresh_token string, err error) {

	claims := &signedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		User_id:    user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refresh_claims := &signedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		User_id:    user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(72)).Unix(),
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

func UpdateAllTokens() {

}

func ValidateAllTokens() {

}
