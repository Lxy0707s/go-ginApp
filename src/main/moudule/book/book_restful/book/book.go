package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/moudule/book/book_service"
	"go-ginApp/src/main/moudule/book/book_service/book"
	mygzip "go-ginApp/src/main/pkg/utils/funcs/mygzip"
	"go-ginApp/src/main/pkg/utils/httptool"
	"go-ginApp/src/main/pkg/utils/myfile"
	"io/ioutil"
	"net/http"
	"time"
)

type Book struct {
	DBId          int32     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DBBookName    string    `gorm:"column:book_name" json:"book_name"`
	DBAuthor      string    `gorm:"column:author" json:"author"`
	DBPrice       int32     `gorm:"column:price" json:"price"`
	DBDescribe    string    `gorm:"column:describe" json:"describe"`
	DBReleaseDate time.Time `gorm:"column:release_date" json:"release_date"`
	DBStatus      int32     `gorm:"column:status" json:"status"`
}

// GetBookList 获取所有机器信息
func GetBookList(c *gin.Context) {
	appG := httptool.Gin{C: c}
	_, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	var args book_service.QueryArgs
	if err := c.ShouldBind(&args); err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, "url err")
		return
	}
	var info, _ = book.GetBookList(c, args)
	bytes, errd := json.Marshal(info)
	if errd != nil {
		return
	}
	compressData, err := mygzip.Compress(bytes)
	if err != nil {
		return
	}
	deCompressData, err := mygzip.DeCompress(compressData)
	if err != nil {
		return
	}
	err = myfile.WriteJSON(string(deCompressData), "compressData.json")
	if err != nil {
		return
	}
	var res = make([]*Book, 0)
	json.Unmarshal(deCompressData, &res)
	if info != nil {
		appG.Response(http.StatusOK, httptool.SUCCESS, res)
	} else {
		appG.Response(http.StatusBadRequest, httptool.SUCCESS, nil)
	}
	return
}

// GetBookByName 根据名称获取书籍信息
func GetBookByName(c *gin.Context) {
	appG := httptool.Gin{C: c}
	by, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	//fmt.Println("body:", string(by))
	var args book_service.QueryArgs
	err1 := json.Unmarshal(by, &args)
	if err1 != nil {
		fmt.Println(err1, "ccccccc")
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, errors.New(err1.Error()))
		return
	}
	/*if err := c.ShouldBind(&args); err != nil && err1 != nil{
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, "url err")
		return
	}*/
	fmt.Println("参数:: ", args.BookName)
	var info, _ = book.GetBookList(c, args)
	bytes, errd := json.Marshal(info)
	if errd != nil {
		return
	}
	compressData, err := mygzip.Compress(bytes)
	if err != nil {
		return
	}
	deCompressData, err := mygzip.DeCompress(compressData)
	if err != nil {
		return
	}
	err = myfile.WriteJSON(string(deCompressData), "compressData.json")
	if err != nil {
		return
	}
	var res = make([]*Book, 0)
	json.Unmarshal(deCompressData, &res)
	if info != nil {
		appG.Response(http.StatusOK, httptool.SUCCESS, res)
	} else {
		appG.Response(http.StatusBadRequest, httptool.SUCCESS, nil)
	}
	return
}

// DeleteBookByID 根据id删除书籍信息
func DeleteBookByID(c *gin.Context) {
	appG := httptool.Gin{C: c}
	by, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, map[string]interface{}{
			"error": errors.New("body readAll error"),
		})
		return
	}
	var args book_service.QueryArgs
	err1 := json.Unmarshal(by, &args)
	if err1 != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, errors.New(err1.Error()))
		return
	}
	fmt.Println("参数:: ", args.Id)
	err = book.DeleteBook(c, args)
	if err != nil {
		appG.Response(http.StatusBadRequest, httptool.InvalidParams, nil)
	}
	appG.Response(http.StatusOK, httptool.SUCCESS, "ok")
	return
}
