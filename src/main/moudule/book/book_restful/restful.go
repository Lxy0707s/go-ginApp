package book_restful

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/moudule/book/book_restful/book"
)

var routeModuleName = "book"

func RegisterRoute(router *gin.RouterGroup) {
	//	注册
	registerRestfulApiV1(router)
}

func registerRestfulApiV1(router *gin.RouterGroup) {
	router = router.Group("/v1/" + routeModuleName)
	//查询接口
	router.GET("/getBookList", book.GetBookList)
	router.POST("/getBookInfoByName", book.GetBookByName)

	router.POST("/deleteBook", book.DeleteBookByID)
}
