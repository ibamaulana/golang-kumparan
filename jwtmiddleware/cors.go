package jwtmiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamaulana/golang-kumparan/config"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conf := config.NewConfig()

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", conf.GetString("cors.allowed_origins"))
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", conf.GetString("cors.allowed_headers"))
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", conf.GetString("cors.allowed_methods"))

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
		}

		ctx.Next()
	}
}
