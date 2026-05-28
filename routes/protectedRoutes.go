package routes

import (
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/movie/:imdb_id", controllers.GetMovie())
	protected.POST("/addMovie", controllers.AddMovie())
	protected.GET("/recommendedMovies", controllers.GetRecommendedMovies())
}
