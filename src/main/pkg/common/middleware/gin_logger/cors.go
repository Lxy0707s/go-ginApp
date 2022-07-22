package gin_logger

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := "*"
		if len(c.Request.Header["Origin"]) > 0 {
			origin = c.Request.Header["Origin"][0]
		} else if len(c.Request.Header["Referer"]) > 0 {
			origin = c.Request.Header["Referer"][0]
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
