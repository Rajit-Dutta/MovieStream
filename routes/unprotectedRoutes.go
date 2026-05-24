package routes

import (
	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetUpUnProtectedRoutes(router *gin.Engine) {
	router.GET("/movies", controllers.GetMovies())
	router.POST("/registerUser", controllers.RegisterUser())
	router.POST("/loginUser", controllers.LoginUser())
}
