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
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movies

		cursor, err := movieCollection.Find(cntxt, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error: Failed to fetch movies"})
			return
		}
		defer cursor.Close(cntxt)

		if err = cursor.All(cntxt, &movies); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error: Failed to decode movies"})
			return
		}

		ctx.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieId := ctx.Param("imdb_id")
		if movieId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "error: Movie ID is required"})
			return
		}

		var movie models.Movies
		err := movieCollection.FindOne(cntxt, bson.M{"imdb_id": movieId}).Decode(&movie)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "error: Failed to find movie"})
			return
		}
		ctx.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cntxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movies

		if err := ctx.ShouldBindJSON(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Validation"})
			return
		}

		result, err := movieCollection.InsertOne(cntxt, movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Body"})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
