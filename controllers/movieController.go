package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/database"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		movieID := ctx.Param("imdb_id")
		if movieID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "error: Movie ID is required"})
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var res struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
			return
		}

		sentiment, rankval, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}

		filter := bson.D{{Key: "imdb_id", Value: movieID}}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankval,
					"ranking_name":  sentiment,
				},
			},
		}

		var cntx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(cntx, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		res.RankingName = sentiment
		res.AdminReview = req.AdminReview

		ctx.JSON(http.StatusOK, res)
	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	if err = godotenv.Load(".env"); err != nil {
		log.Println("ENV file missing!")
	}

	GROQ_api_key := os.Getenv("GROQ_API_KEY")
	if GROQ_api_key == "" {
		return "", 0, errors.New("API KEY missing")
	}

	llm, err := openai.New(openai.WithToken(GROQ_api_key))
	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(context.Background(), base_prompt+admin_review)
	if err != nil {
		return "", 0, err
	}
	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetUsersFavouriteGenres(userId string) ([]string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userId}}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve genres of the user")
	}

	var genreNames []string

	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}
