package dao

import (
	"github.com/jinzhu/gorm"
	"go-ginApp/src/main/pkg/utils/dbtool"
)

// DemoDao 数据库配置层
func DemoDao() *gorm.DB {
	return dbtool.DBMap["demo"] // 初始化数据库名
}
