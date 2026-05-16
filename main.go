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

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server is not yet hosted in :8080")
	}
}
