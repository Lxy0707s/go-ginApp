package httpServer

import (
	"fmt"
	"go-ginApp/src/main/pkg/utils/funcs/hash"
	"go-ginApp/src/main/pkg/utils/funcs/ips"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	hashToken  = "goping"
	hashFormat = "%s-" + hashToken + "-%d"

	keyHeader  = "Hash-Key"
	codeHeader = "Hash-Code"

	hashDuration int64 = 1800
)

// BuildHashHeader build header with dataType and local ip
func BuildHashHeader(dataType string) map[string]string {
	localIP := ips.LocalIP()
	hour := time.Now().Unix() / hashDuration
	key := hash.MD5String(dataType + hashToken + localIP)
	code := hash.MD5String(fmt.Sprintf(hashFormat, key, hour))
	return map[string]string{
		keyHeader:  key,
		codeHeader: code,
	}
}

// CheckHashHandler is a HTTP handler to check hash
func CheckHashHandler(rw http.ResponseWriter, r *http.Request) {
	if !checkHashHeader(r.Header) {
		rw.Header().Add("Connection", "close")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

// CheckHashGinHandler is a gin handler to check hash
func CheckHashGinHandler(c *gin.Context) {
	CheckHashHandler(c.Writer, c.Request)
	if c.Writer.Status() == http.StatusUnauthorized {
		c.Writer.Header().Add("Connection", "close")
		c.Abort()
	}
}

func checkHashHeader(header http.Header) bool {
	key := header.Get(keyHeader)
	if key == "" {
		return false
	}
	code := header.Get(codeHeader)
	if code == "" {
		return false
	}
	hour := time.Now().Unix() / hashDuration
	myCode := hash.MD5String(fmt.Sprintf(hashFormat, key, hour))
	if myCode == code {
		return true
	}
	hourBefore := hour - 1
	myCode = hash.MD5String(fmt.Sprintf(hashFormat, key, hourBefore))
	if myCode == code {
		return true
	}

	hourAfter := hour + 1
	myCode = hash.MD5String(fmt.Sprintf(hashFormat, key, hourAfter))
	return myCode == code
}
