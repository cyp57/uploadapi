package response

import (
	"github.com/cyp57/uploadapi/cnst"
	"github.com/gin-gonic/gin"
)

type IResponseHandler interface {
	SuccessResponse(*gin.Context, int, interface{}, string)
	ErrResponse(*gin.Context, int, string)
	DataResponse(*gin.Context, int, interface{})
}

type (
	ResponseHandler struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Status  bool        `json:"status"`
		Error   *ErrorRes   `json:"error"`
	}

	ErrorRes struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
	}
)

func Response() IResponseHandler {
	return &ResponseHandler{}
}

var ErrCode = 0

func SetErrorCode(errCode int) {
	ErrCode = errCode
}

func (r *ResponseHandler) SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, &ResponseHandler{Data: data, Message: message, Status: true})
}

func (r *ResponseHandler) ErrResponse(c *gin.Context, statusCode int, message string) {
	if ErrCode != 0 && message == "" {
		errMsg := cnst.GetErrMsg(ErrCode)
		message = errMsg
	}
	c.AbortWithStatusJSON(statusCode, &ResponseHandler{Status: false, Error: &ErrorRes{ErrMsg: message, ErrCode: ErrCode}})

}

func (r *ResponseHandler) DataResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}
