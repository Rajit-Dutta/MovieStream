package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/database"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

var userCollection *mongo.Collection = database.OpenCollection("users")

func RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		validate := validator.New()
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Problem in hashing password"})
			return
		}

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Problem in getting user input"})
			return
		}
		if err := validate.Struct(user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
			return
		}
		count, err := userCollection.CountDocuments(cntxt, bson.M{"email": user.Email})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error in retrieving user duplicate data"})
			return
		}
		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(cntxt, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in registering user"})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
