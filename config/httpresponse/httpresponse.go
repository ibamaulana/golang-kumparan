package httpresponse

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ResponseException struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
}

type ResponsePaged struct {
	Data interface{} `json:"data"`
	Page int         `json:"page"`
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

func NewSuccessResponsePaged(c *gin.Context, data interface{}, page int) {
	response := new(ResponsePaged)
	response.Data = data
	response.Page = page
	c.JSON(http.StatusOK, &response)
	return
}

func flagErr() string {
	return "[ERROR-" + time.Now().Format("20060102-150405.0000") + "]"
}
