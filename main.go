package main

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/ibamaulana/golang-kumparan/config/jwtmiddleware"
	"github.com/ibamaulana/golang-kumparan/controller/news"
	"github.com/ibamaulana/golang-kumparan/controller/ping"
	"github.com/ibamaulana/golang-kumparan/queue"
)

func main() {

	queue.Consumestart()
	r := gin.Default()

	signInKey := "secret"
	jwtmiddleware.InitJWTMiddlewareCustom([]byte(signInKey), jwt.SigningMethodHS512)

	r.Use(jwtmiddleware.CORSMiddleware())
	r.GET("ping", ping.PingController)

	{
		newsRoute := r.Group("news")
		// newsRoute.Use(jwtmiddleware.MyAuth())

		newsRoute.POST("", news.CreateController)
		newsRoute.GET("", news.GetController)
		newsRoute.GET("cache", news.CacheController)
	}

	r.Run(":9000")

}
