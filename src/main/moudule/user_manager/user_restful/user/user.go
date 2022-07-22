package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/moudule/user_manager/user_service"
	"go-ginApp/src/main/moudule/user_manager/user_service/user"
	"go-ginApp/src/main/pkg/utils/httptool"
	"io/ioutil"
	"net/http"
)

// GetUserByEmail 根据email获取用户信息
func GetUserByEmail(c *gin.Context) {
	appG := httptool.Gin{C: c}
	by, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	//fmt.Println("body:", string(by))
	var args user_service.QueryArgs
	err1 := json.Unmarshal(by, &args)
	if err1 != nil {
		var info, _ = user.GetUserList(c, args)
		appG.Response(http.StatusOK, httptool.InvalidParams, info)
		return
	}
	fmt.Println("user参数:: ", args.UserName)
	var info, _ = user.GetUserList(c, args)
	if info != nil {
		appG.Response(http.StatusOK, httptool.SUCCESS, info)
	} else {
		appG.Response(http.StatusBadRequest, httptool.SUCCESS, nil)
	}
	return
}

// RegisterUser 注册
func RegisterUser(c *gin.Context) {
	appG := httptool.Gin{C: c}
	by, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	//fmt.Println("body:", string(by))
	var args user_service.QueryArgs
	err1 := json.Unmarshal(by, &args)
	if err1 != nil {
		appG.Response(http.StatusOK, httptool.InvalidParams, nil)
		return
	}
	res, _ := user.QueryUserInfo(c, "")
	appG.Response(http.StatusOK, httptool.InvalidParams, res)
	return
}
