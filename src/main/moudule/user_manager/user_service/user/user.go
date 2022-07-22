package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/dao"
	db "go-ginApp/src/main/internal/dao/db_models"
	"go-ginApp/src/main/moudule/book/book_service"
	"go-ginApp/src/main/moudule/user_manager/user_service"
)

func QueryUserInfo(ctx *gin.Context, useName string) (*db.UserDB, error) {
	var bookList db.UserDB
	//	数据库连接
	dbResult := dao.DemoDao().Model(db.UserDB{})
	//	预加载
	if useName != "" {
		dbResult = dbResult.Where("user_name = ?", useName)
	}
	dbResult = dbResult.Order("id DESC").First(&bookList)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &bookList, nil
}

// GetUserList restful接口专用
func GetUserList(ctx *gin.Context, queryArgs user_service.QueryArgs) (*[]*db.UserDB, error) {
	var bookList []*db.UserDB
	//	数据库连接
	dbResult := dao.DemoDao().Model(db.UserDB{})
	//	预加载
	fmt.Println(queryArgs.UserName, "----------")
	if &queryArgs.UserName != nil && queryArgs.UserName != "" {
		dbResult = dbResult.Where("user_name = ?", queryArgs.UserName)
	}
	dbResult = dbResult.Order("id DESC").Scan(&bookList)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &bookList, nil
}

// DeleteUser restful接口专用
func DeleteUser(ctx *gin.Context, queryArgs book_service.QueryArgs) error {

	return nil
}
