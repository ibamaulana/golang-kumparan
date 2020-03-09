package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamaulana/golang-master/config"
)

func PingController(ctx *gin.Context) {
	cfg := config.NewConfig()

	db, err := config.MysqlConnection(cfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	_ = db

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	return
}
