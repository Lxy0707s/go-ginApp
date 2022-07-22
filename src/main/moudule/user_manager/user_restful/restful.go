package user_restful

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/moudule/user_manager/user_restful/user"
)

var routeModuleName = "user"

func RegisterRoute(router *gin.RouterGroup) {
	//	注册
	registerRestfulApiV1(router)
}

func registerRestfulApiV1(router *gin.RouterGroup) {
	router = router.Group("/v1/" + routeModuleName)
	//创建用户接口
	router.OPTIONS("/register", user.RegisterUser)
	//更新用户接口
	//router.GET("/updateUser", user.GetUserByEmail)
	//注销用户接口
	//router.GET("/deleteUser", user.GetUserByEmail)
	//查询用户接口
	router.GET("/queryUser", user.GetUserByEmail)
}
