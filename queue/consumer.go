package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adjust/rmq"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gin-gonic/gin"
	"github.com/ibamaulana/golang-kumparan/config"
	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/ibamaulana/golang-kumparan/services"
)

func Consumestart() {
	cfg := config.NewConfig()
	redisport := cfg.GetString(`redis.port`)
	connectionconsumer := rmq.OpenConnection("consumer", "tcp", redisport, 2)

	newstask := connectionconsumer.OpenQueue("newsTask")
	newstask.StartConsuming(10, time.Second)
	newstask.AddConsumer("consumer", NewConsumer(1))
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

	//parsing payload
	if err = json.Unmarshal([]byte(delivery.Payload()), &news); err != nil {
		// handle error
		delivery.Reject()
		return
	}

	// perform task
	log.Printf("performing task %s", news.Author)

	//save to database
	cfg := config.NewConfig()
	db, err := config.MysqlConnection(cfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//get service
	newsContract := services.NewNewsServiceContract(db)
	tx := db.Begin()
	defer tx.Rollback()

	//run create on db
	news, err = newsContract.Create(news, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tx.Commit()

	//run create on elatic search
	elasticrequest := model.ElasticRequest{
		Created: news.Created,
	}

	newsdata, err := json.Marshal(elasticrequest)
	if err != nil {
		return
	}

	//Set up the request object.
	req := esapi.IndexRequest{
		Index:      "news",
		DocumentID: strconv.Itoa(news.ID),
		Body:       bytes.NewReader(newsdata),
		Refresh:    "true",
	}

	elasticcfg := config.ElasticConnection(cfg)

	//creating new client
	es, err := elasticsearch.NewClient(elasticcfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), news.ID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	delivery.Ack()
}
