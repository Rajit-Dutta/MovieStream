package main

import (
	"log"
	"net/http"

	"github.com/Rajit-Dutta/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	routes.SetUpProtectedRoutes(router)
	routes.SetUpUnProtectedRoutes(router)

	if err := router.Run(":8082"); err != nil {
		log.Fatal("Server is not yet hosted in :8080")
	}
}
