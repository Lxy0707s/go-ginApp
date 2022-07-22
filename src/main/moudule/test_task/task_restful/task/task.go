package task

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	db "go-ginApp/src/main/internal/dao/db_models"
	"go-ginApp/src/main/moudule/test_task/task_service"
	"go-ginApp/src/main/moudule/test_task/task_service/task"
	"go-ginApp/src/main/pkg/utils/gziptool"
	"go-ginApp/src/main/pkg/utils/httptool"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	gzipWriter *gzip.Writer
	err        error
)

// GetAll 获取所有机器信息
func GetAll(c *gin.Context) {
	appG := httptool.Gin{C: c}
	_, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	var args task_service.QueryArgs
	if err := c.ShouldBind(&args); err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, "url err")
		return
	}
	var ss = make([]*db.DemoTask, 0)
	var info, _ = task.GetTaskInfo(c, args)
	// 压缩流程
	var buf bytes.Buffer
	jsonByte, _ := json.Marshal(info)
	err = gziptool.GzipWrite(&buf, jsonByte)
	if err != nil {
		log.Fatal(err)
	}
	// 解压流程
	var buf2 bytes.Buffer
	err = gziptool.GunzipWrite(&buf2, buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(buf2.Bytes(), &ss)
	if err != nil {
		log.Fatal(err)
		return
	}
	if info != nil {
		appG.Response(http.StatusOK, httptool.SUCCESS, buf.Bytes())
	} else {
		appG.Response(http.StatusBadRequest, httptool.SUCCESS, nil)
	}
	return
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, response{
		Code: 200,
		Msg:  "success",
		Data: "hello",
	})
	return
}
