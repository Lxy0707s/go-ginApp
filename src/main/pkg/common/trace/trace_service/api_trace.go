package trace_service

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"go-ginApp/src/main/pkg/common/trace/models"
	"log"
)

func ApiLogTrace(ctx context.Context, tableName string, logTypeDB int, data interface{}) error {
	// 只有设置了ctx才需进行记录，避免计划任务大量记录
	if ctx == nil {
		return nil
	}
	// 需要来自接口请求
	logAPIRequestID := GetLogAPIRequestID(ctx)
	if logAPIRequestID == "" {
		return nil
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
	}
	// 开始判断是否需要记录到db
	logMysqlSwitch, logMysqlORM := GetLogMysql(ctx)
	if logMysqlSwitch && logMysqlORM != nil {
		if err := logMysqlORM.Create(&models.OperateLogDetail{
			RequestId:     logAPIRequestID,
			OpTableName:   tableName,
			TableAction:   int32(logTypeDB),
			ChangeContent: string(jsonBytes),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetLogMysql(ctx context.Context) (bool, *gorm.DB) {
	if ctx != nil {
		if ctx.Value("logMysqlSwitch") != nil && ctx.Value("logMysqlORM") != nil {
			return ctx.Value("logMysqlSwitch").(bool), ctx.Value("logMysqlORM").(*gorm.DB)
		}
	}

	return false, nil
}

func GetIp(ctx context.Context) string {
	if ctx != nil {
		if ctx.Value("logAPIRequestIP") != nil {
			return ctx.Value("logAPIRequestIP").(string)
		}
	}

	return ""
}

func GetLogAPIRequestID(ctx context.Context) string {
	if ctx != nil {
		if ctx.Value("logAPIRequestID") != nil {
			return ctx.Value("logAPIRequestID").(string)
		}
	}

	return ""
}
