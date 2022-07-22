package task_restful

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/moudule/test_task/task_restful/task"
)

var routeModuleName = "test"

func RegisterRoute(router *gin.RouterGroup) {
	//	注册
	registerRestfulApiV1(router)
}

func registerRestfulApiV1(router *gin.RouterGroup) {
	router = router.Group("/v1/" + routeModuleName)
	//查询接口
	router.GET("/getAllInfo", task.GetAll)
	router.GET("Hello", task.Hello)
}
