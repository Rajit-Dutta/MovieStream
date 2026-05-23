package main

import (
	"log"
	"net/http"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movie/:imdb_id", controllers.GetMovie())
	router.POST("/addMovie", controllers.AddMovie())
	router.POST("/registerUser", controllers.RegisterUser())
	router.POST("/loginUser", controllers.LoginUser())

	if err := router.Run(":8082"); err != nil {
		log.Fatal("Server is not yet hosted in :8080")
	}
}
