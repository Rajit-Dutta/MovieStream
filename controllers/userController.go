package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/database"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/models"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/utils"
	"github.com/gin-gonic/gin"
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

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Problem in getting user input"})
			return
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Problem in hashing password"})
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

func LoginUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var loginUser models.UserLogin
		var foundUser models.User

		if err := ctx.ShouldBindJSON(&loginUser); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Problem in getting login user input"})
			return
		}

		if err := validate.Struct(loginUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
			return
		}

		err := userCollection.FindOne(cntxt, bson.M{"email": loginUser.Email}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginUser.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Password mismatch"})
			return
		}

		accessToken, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateALlTokens(foundUser.UserID, accessToken, refreshToken)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		ctx.JSON(http.StatusOK, models.UserResponse{
			UserId:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			Token:           accessToken,
			RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}
