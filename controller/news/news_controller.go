package news

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/adjust/rmq"
	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ibamaulana/golang-kumparan/config"
	"github.com/ibamaulana/golang-kumparan/config/cache"
	"github.com/ibamaulana/golang-kumparan/config/httpresponse"
	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/ibamaulana/golang-kumparan/request/news"
	"github.com/ibamaulana/golang-kumparan/services"
	"github.com/jinzhu/copier"
)

func CreateController(ctx *gin.Context) {
	runtime.GOMAXPROCS(1)
	cfg := config.NewConfig()
	var err error

	var req news.CreateRequest

	if err = ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	news := new(model.News)
	err = copier.Copy(&news, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//produce
	redisport := cfg.GetString(`redis.port`)
	connection := rmq.OpenConnection("producer", "tcp", redisport, 2)
	newsTask := connection.OpenQueue("newsTask")
	newsTaskData, err := json.Marshal(news)
	if err != nil {
		return
	}
	newsTask.PublishBytes(newsTaskData)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	return
}

func GetController(ctx *gin.Context) {
	runtime.GOMAXPROCS(2)
	var m map[string]interface{}
	var data []*model.ElasticResponse
	cfg := config.NewConfig()

	//fetching data
	go FetchController()

	//get data from elastic
	elasticcfg := config.ElasticConnection(cfg)
	es, err := elasticsearch.NewClient(elasticcfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	//pagination
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	start := (page - 1) * 10

	//elastic query sort by created
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": map[string]interface{}{
			"created": map[string]interface{}{
				"order": "desc",
			},
		},
		"from": start, "size": 10,
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("news"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// parsing from elastic
	for _, hit := range m["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, _ := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))
		date, _ := time.Parse(time.RFC3339, hit.(map[string]interface{})["_source"].(map[string]interface{})["created"].(string))

		datas := new(model.ElasticResponse)
		datas.ID = id
		datas.Created = date
		data = append(data, datas)
	}

	httpresponse.NewSuccessResponsePaged(ctx, data, page)
	return
}

func FetchController() {
	var err error
	var ctx *gin.Context
	cfg := config.NewConfig()

	db, err := config.MysqlConnection(cfg)
	if err != nil {
		httpresponse.NewErrorException(ctx, http.StatusBadRequest, err)
		return
	}

	newsContract := services.NewNewsServiceContract(db)
	data, err := newsContract.Get()

	cache.SetCache("datanews", data)
}

func CacheController(ctx *gin.Context) {
	data, found := cache.GetCache("datanews")
	if found {
		httpresponse.NewSuccessResponse(ctx, data)
		return
	}

	httpresponse.NewSuccessResponse(ctx, nil)
	return
}
