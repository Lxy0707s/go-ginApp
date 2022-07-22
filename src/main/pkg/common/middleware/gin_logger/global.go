package gin_logger

import (
	"github.com/gin-gonic/gin"
)

// Global middleware
func Global(r *gin.Engine) *gin.Engine {
	//r.Use(gin.Logger())
	r.Use(Logger())
	r.Use(CORS())
	return r
}
