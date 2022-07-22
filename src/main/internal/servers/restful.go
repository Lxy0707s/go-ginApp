package servers

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/middleware"
	"go-ginApp/src/main/moudule/book/book_restful"
	"go-ginApp/src/main/moudule/test_task/task_restful"
	"go-ginApp/src/main/moudule/user_manager/user_restful"
)

func RegisterRestfulRoute(r *gin.Engine, s *HTTPServer) {
	//gin 路由初始化
	apiG := r.Group("/api")
	//sso路由注册
	if true { // debug 开启时才允许访问iql
		/*r.GET("/graphiql", gin.WrapF(graphiql.NewGraphiQLHandlerFunc()))
		r.GET("/playground", playgroundHandler())*/
	}
	apiG.Use(s.profilerMiddleware)
	//apiG.Use(sys_jwt.ApiTokenAuth(config.AppConfig.ApiTokens)) // restful jwt校验中间件
	apiG.Use(middleware.NewSysInstance().UseJwtCheck())
	book_restful.RegisterRoute(apiG)
	task_restful.RegisterRoute(apiG)
	user_restful.RegisterRoute(apiG)
}

/*func playgroundHandler() gin.HandlerFunc {
	h := graphiql.Playground("GraphQL", "/graphql/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}*/
