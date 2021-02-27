package main

import (
	"devfelipereis/urlShortener/env"
	"devfelipereis/urlShortener/routes"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	appEnv := env.Get()
	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	router.GET("/", routes.Home)
	router.GET("/:code", routes.Redirect)
	router.GET("/:code/info", routes.GetOne)
	router.POST("/generate", routes.Generate)

	router.Run(appEnv.ApiPort)
}
