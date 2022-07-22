package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/config"
	"go-ginApp/src/main/moudule/user_manager/user_service/user"
	"go-ginApp/src/main/pkg/common/middleware/sys_jwt"
	"go-ginApp/src/main/pkg/utils/httptool"
	"net/http"
	"strings"
	"time"
)

type (
	SysOp struct {
		jwt sys_jwt.JwtImpl
	}
)

func NewSysInstance() *SysOp {
	return &SysOp{
		jwt: sys_jwt.NewJwtInstance(),
	}
}

func (s *SysOp) UseJwtCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := httptool.Gin{C: c}
		if c.Request.Method != "OPTIONS" {
			token := c.GetHeader("Authorization")
			userName := c.GetHeader("user_name")
			token = strings.Replace(token, "Bearer ", "", 1)
			if token == "" {
				token = c.Query("jwt_token")
			}
			if token == "" {
				var info, _ = user.QueryUserInfo(c, userName) // 这里可以使用任何形式的用户验证接口，同步远程验证结果，进行下一步处理
				if info == nil {
					c.Abort()
					return
				}
				userName = info.DBUserName
				token = info.DBToken
			}
			if userName == "" {
				userName = "test"
			}
			fmt.Println(userName)
			check, err := s.jwt.ParseJwtToken(token, userName)
			if err != nil {
				config.AppLog.Warn("error:", "jwt check error:", err)
				appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
					"result": err.Error(),
				})
				c.Abort()
				return
			}
			c.Set("zName", check.Uname)
			c.Set("zEmail", check.Email)
		}
		c.Next()
	}
}

// 需要校验的字段，返回校验结果和报错信息
func (s *SysOp) jwtCheck(params sys_jwt.SysClaims) (bool, error) {
	var jwtParams, err = json.Marshal(params)
	var jwtClaims *sys_jwt.SysClaims
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(jwtParams, &jwtClaims)
	if err != nil || jwtClaims == nil {
		return false, err
	}
	//检查一下超时没有,超时就返回错误
	if jwtClaims.VerifyExpiresAt(int64(time.Now().Unix()), true) == false {
		return false, errors.New("超时")
	}
	//检查一下发行方正确,令牌的发行方错误
	if !jwtClaims.VerifyIssuer(AppIss, true) {
		return false, errors.New("token's issuer is wrong")
	}
	return true, nil
}

func GetUserInfo(ctx context.Context) (string, string) {
	var user = ""
	var cname = ""
	if ctx != nil {
		if ctx.Value("zName") != nil {
			user = ctx.Value("zName").(string)
		}
		if ctx.Value("zEmail") != nil {
			cname = ctx.Value("zEmail").(string)
		}
	}

	return user, cname
}
