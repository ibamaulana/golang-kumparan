package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ibamaulana/golang-kumparan/config"
	"github.com/ibamaulana/golang-kumparan/controller/news"
	"github.com/ibamaulana/golang-kumparan/controller/ping"
	"github.com/ibamaulana/golang-kumparan/jwtmiddleware"
	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/ibamaulana/golang-kumparan/services"
)

func main() {
	connectionconsumer := rmq.OpenConnection("consumer", "tcp", "localhost:6379", 2)

	newstask := connectionconsumer.OpenQueue("newsTask")
	newstask.StartConsuming(10, time.Second)
	newstask.AddConsumer("consumer", NewConsumer(1))

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
		newsRoute.GET(":id", news.FindController)
	}

	r.Run(":9000")

}

type Consumer struct {
	name   string
	count  int
	before time.Time
}

func NewConsumer(tag int) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("consumer%d", tag),
		count:  0,
		before: time.Now(),
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	var err error
	var news *model.News
	var ctx *gin.Context

	fmt.Println(string(delivery.Payload()))
	if err = json.Unmarshal([]byte(delivery.Payload()), &news); err != nil {
		// handle error
		delivery.Reject()
		return
	}

	// perform task
	log.Printf("performing task %s", news.Author)
	cfg := config.NewConfig()
	db, err := config.MysqlConnection(cfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newsContract := services.NewNewsServiceContract(db)
	tx := db.Begin()
	defer tx.Rollback()

	err = newsContract.Create(news, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tx.Commit()

	delivery.Ack()
}
