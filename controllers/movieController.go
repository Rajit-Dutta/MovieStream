package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/database"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func GetMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movies

		cursor, err := movieCollection.Find(cntxt, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error: Failed to fetch movies"})
		}
		defer cursor.Close(cntxt)

		if err = cursor.All(cntxt, &movies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error: Failed to decode movies"})
		}

		ctx.JSON(http.StatusOK, movies)
	}
}
