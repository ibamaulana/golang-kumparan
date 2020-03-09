package news

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

	"github.com/adjust/rmq"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ibamaulana/golang-kumparan/config"
	"github.com/ibamaulana/golang-kumparan/httpresponse"
	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/ibamaulana/golang-kumparan/request/news"
	"github.com/ibamaulana/golang-kumparan/services"
	"github.com/jinzhu/copier"
)

func CreateController(ctx *gin.Context) {
	runtime.GOMAXPROCS(1)

	var err error
	cfg := config.NewConfig()

	db, err := config.MysqlConnection(cfg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var req news.CreateRequest

	if err = ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newsContract := services.NewNewsServiceContract(db)

	tx := db.Begin()
	defer tx.Rollback()

	news := new(model.News)
	err = copier.Copy(&news, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//produce
	connection := rmq.OpenConnection("producer", "tcp", "localhost:6379", 2)
	newsTask := connection.OpenQueue("newsTask")
	newsTaskData, err := json.Marshal(news)
	if err != nil {
		return
	}
	newsTask.PublishBytes(newsTaskData)

	err = newsContract.Create(news, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

	return
}

func GetController(ctx *gin.Context) {
	runtime.GOMAXPROCS(2)

	var err error
	cfg := config.NewConfig()

	db, err := config.MysqlConnection(cfg)
	if err != nil {
		httpresponse.NewErrorException(ctx, http.StatusBadRequest, err)
		return
	}

	// newsContract := services.NewNewsServiceContract(db)
	// data, err := newsContract.Get()
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "3"))
	var news []*model.News

	paginator := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"created desc"},
		ShowSQL: true,
	}, &news)

	httpresponse.NewSuccessResponse(ctx, paginator)
	return
}

func FindController(ctx *gin.Context) {
	runtime.GOMAXPROCS(2)

	var err error
	cfg := config.NewConfig()

	db, err := config.MysqlConnection(cfg)
	if err != nil {
		httpresponse.NewErrorException(ctx, http.StatusBadRequest, err)
		return
	}

	var req news.FindRequest

	if err = ctx.ShouldBindUri(&req); err != nil {
		httpresponse.NewErrorException(ctx, http.StatusBadRequest, err)
		return
	}

	newsContract := services.NewNewsServiceContract(db)
	data, err := newsContract.Find(req.ID)

	httpresponse.NewSuccessResponse(ctx, data)
	return
}
