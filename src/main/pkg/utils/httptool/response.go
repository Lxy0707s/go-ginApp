package httptool

import (
	"github.com/gin-gonic/gin"
)

// Gin Gin 返回的包装
type Gin struct {
	C *gin.Context
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//Response 包装后的信息，带业务代码和消息
func (g *Gin) Response(httpCode, code int, data interface{}, freeMsg ...string) {
	var msg string
	if len(freeMsg) > 0 && freeMsg[0] != "" {
		msg = freeMsg[0]
	} else {
		msg = GetMsg(code)
	}
	g.C.JSON(httpCode, response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	return
}

// DataResponse 直接返回 json
func (g *Gin) DataResponse(httpCode int, data interface{}) {
	g.C.JSON(httpCode, data)
	return
}
