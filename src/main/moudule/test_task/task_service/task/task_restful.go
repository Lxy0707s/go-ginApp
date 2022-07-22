package task

import (
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/dao"
	db "go-ginApp/src/main/internal/dao/db_models"
	"go-ginApp/src/main/moudule/test_task/task_service"
)

// GetTaskInfo restful接口专用
func GetTaskInfo(ctx *gin.Context, queryArgs task_service.QueryArgs) (*[]*db.DemoTask, error) {
	var taskList []*db.DemoTask
	//	数据库连接
	dbResult := dao.DemoDao()
	//	预加载
	if &queryArgs.TaskName != nil && queryArgs.TaskName != "" {
		dbResult = dbResult.Where("net_detect_new.task_switch != ?", 2)
	}
	if queryArgs.TaskId != "" {
		dbResult = dbResult.Where("net_detect_new.task_switch != ?", 2)
	}
	dbResult = dbResult.Model(db.DemoTask{}).Order("id DESC").Find(&taskList)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &taskList, nil
}
