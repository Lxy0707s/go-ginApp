package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"go-ginApp/src/main/moudule/graphql_demo/graphql_query"
)

func GraphqlHandler() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   &graphql_query.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// 只需要通过Gin简单封装即可
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
