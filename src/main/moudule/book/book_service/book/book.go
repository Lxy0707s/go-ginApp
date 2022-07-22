package book

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/dao"
	db "go-ginApp/src/main/internal/dao/db_models"
	"go-ginApp/src/main/moudule/book/book_service"
)

// GetBookList restful接口专用
func GetBookList(ctx *gin.Context, queryArgs book_service.QueryArgs) (*[]*db.BookDB, error) {
	var bookList []*db.BookDB
	//	数据库连接
	dbResult := dao.DemoDao().Model(db.BookDB{})
	//	预加载
	fmt.Println(queryArgs.BookName, "----------")
	if &queryArgs.BookName != nil && queryArgs.BookName != "" {
		dbResult = dbResult.Where("book_name = ?", queryArgs.BookName)
	}
	dbResult = dbResult.Order("id DESC").Scan(&bookList)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &bookList, nil
}

// DeleteBook restful接口专用
func DeleteBook(ctx *gin.Context, queryArgs book_service.QueryArgs) error {
	var bookList []*db.BookDB
	//	数据库连接
	dbResult := dao.DemoDao().Model(db.BookDB{})
	//	预加载
	fmt.Println(queryArgs.Id, "----------")
	if &queryArgs.Id != nil && queryArgs.Id != "" {
		dbResult = dbResult.Where("id = ?", queryArgs.Id)
	}
	dbResult = dbResult.Delete(&bookList)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}
