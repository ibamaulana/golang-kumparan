package httpresponse

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type ResponseException struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
}

type ResponsePaged struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
}

type ResponseObject struct {
	Data interface{} `json:"data"`
}

func NewErrorException(c *gin.Context, code int, err error) {
	log.Println(flagErr(), err)

	response := new(ResponseException)
	response.Status = code
	response.Message = err.Error()
	response.ErrorCode = code
	c.JSON(code, &response)

	return
}

func NewSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
	return
}

func flagErr() string {
	return "[ERROR-" + time.Now().Format("20060102-150405.0000") + "]"
}
