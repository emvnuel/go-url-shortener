package main

import (
	"math/rand"
	"time"

	"github.com/emvunel/urlshortener/controllers"
	"github.com/emvunel/urlshortener/models"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("web/*")

	models.ConnectDatabase()

	rand.Seed(time.Now().UnixNano())

	r.POST("/shorten-url", controllers.ShortenUrl)
	r.GET("/:id", controllers.RedirectUrl)

	r.Run()
}
