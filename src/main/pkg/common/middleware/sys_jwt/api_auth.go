package sys_jwt

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/pkg/utils/base_struct"
	e "go-ginApp/src/main/pkg/utils/httptool"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var reqRuleLimit RuleLimit
var baseApiTokenMap map[string]base_struct.ApiTokenMap

type RuleLimit struct {
	lastReqLimitTime    int64
	queryCountLimitReal int
	queryLimitLock      sync.RWMutex
}

type SchedulerParam struct {
	Config struct {
		Key         string `json:"key"`
		Env         string `json:"env"`
		Version     string `json:"version"`
		Type        string `json:"type"`
		Token       string `json:"token"`
		ConfigValue string `json:"configValue"`
		Comments    string `json:"comments"`
	} `json:"config"`
	System   string `json:"system"`
	AppName  string `json:"appName"`
	AppToken string `json:"appToken"`
}

func ApiTokenAuth(apiTokenMap map[string]base_struct.ApiTokenMap) gin.HandlerFunc {
	baseApiTokenMap = apiTokenMap
	return func(c *gin.Context) {
		result := make(map[string]interface{})
		// 支持鉴权请求头中的Authorization[优先]或uri中的appToken
		appToken := FetchToken(c)

		appG := e.Gin{C: c}
		//校验token不可为空
		if appToken == "" {
			appG.Response(http.StatusForbidden, e.EmptyAPIAuth, result)
			c.Abort()
			return
		}

		//校验token是否在配置名单中
		if _, ok := apiTokenMap[appToken]; !ok {
			appG.Response(http.StatusForbidden, e.ErrorAPIAuth, result)
			c.Abort()
			return
		}
		//必备配置参数不可缺少
		if item, ok := apiTokenMap[appToken]; ok {
			if item.AppID == "" ||
				item.AppToken == "" ||
				item.AimApp == "" ||
				item.SourceApp == "" {
				appG.Response(http.StatusUnauthorized, e.ErrorAPIConfig, result)
				c.Abort()
				return
			}

			//根据token匹配调用限制规则
			if item.QueryCountLimit != 0 && item.QueryTimeLimit != 0 {
				reqRuleLimit.queryLimitLock.Lock()
				//校验请求时间和请求次数是否刷新重置
				if reqRuleLimit.lastReqLimitTime != 0 &&
					time.Now().Unix()-int64(item.QueryTimeLimit) > reqRuleLimit.lastReqLimitTime {
					reqRuleLimit.lastReqLimitTime = 0
					reqRuleLimit.queryCountLimitReal = 0
				}
				//每次请求进行信息累计
				reqRuleLimit.queryCountLimitReal += 1
				if reqRuleLimit.lastReqLimitTime == 0 {
					reqRuleLimit.lastReqLimitTime = time.Now().Unix()
				}
				reqRuleLimit.queryLimitLock.Unlock()
				//校验请求次数是否到达规定限制
				if reqRuleLimit.queryCountLimitReal > apiTokenMap[appToken].QueryCountLimit {
					appG.Response(http.StatusForbidden, e.ErrorAPILimit, result)
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}

func FetchSourceApp(c *gin.Context) string {
	token := FetchToken(c)
	if single, ok := baseApiTokenMap[token]; ok {
		return single.SourceApp
	}
	return ""
}

func FetchToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("appToken")
	}
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var params SchedulerParam
	if err := json.Unmarshal(bodyBytes, &params); err == nil && params.AppToken != "" {
		return params.AppToken
	}
	return token
}
