package httptool

import (
	"github.com/gin-gonic/gin"
)

func GetConsumerInfo(c *gin.Context) map[string]string {
	ConsumerMap := make(map[string]string)
	ConsumerMap["id"] = c.GetHeader("X-Consumer-ID")
	ConsumerMap["username"] = c.GetHeader("X-Consumer-Username")
	return ConsumerMap
}
