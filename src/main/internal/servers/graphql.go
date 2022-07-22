package servers

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/graphql/schema"
	"go-ginApp/src/main/internal/middleware"
)

func RegisterGraphqlRoute(r *gin.Engine) {
	// GET方法用于支持GraphQL的web界面操作
	// 如果不需要web界面，可根据自己需求用GET/POST或者其他都可以
	r.Use(middleware.NewSysInstance().UseJwtCheck())
	r.POST("/graphql", schema.GraphqlHandler())
	r.GET("/graphql", schema.GraphqlHandler())
}
