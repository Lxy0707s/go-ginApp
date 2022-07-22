package dbtool

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // gorm 的 MySQL 驱动
)

var (
	// DBMap 数据库多实例
	DBMap map[string]*gorm.DB
	// 数据库配置
	MysqlStore []Option
)

// Option 数据库配置
type Option struct {
	Driver         string
	Username       string
	Password       string
	Addr           string
	DBNames        []string // 多数据库
	DatabasePrefix string
}

// Model 基本类型，自带添加时间和更新时间
type Model struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Extra 增加外部数据，标记不要 db 查询，json 渲染
type Extra struct {
	Ctx context.Context `json:"-" gorm:"-"`
}

// Setup 注册数据库，如果有问题会 panic
func Setup(option interface{}, debug bool) map[string]*gorm.DB {
	switch option.(type) {
	case Option:
		DBMap = initDbConfig(option.(Option), debug)
	case []Option:
		MysqlStore = option.([]Option)
		for _, dbConfig := range option.([]Option) {
			DBMap = initDbConfig(dbConfig, debug)
		}
	}
	return DBMap
}

func initDbConfig(option Option, debug bool) map[string]*gorm.DB {
	if option.Addr == "" {
		return nil
	}
	if DBMap == nil {
		DBMap = make(map[string]*gorm.DB)
	}
	var err error
	for _, name := range option.DBNames {
		key := option.DatabasePrefix + name
		DBMap[key], err = gorm.Open(option.Driver, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&maxAllowedPacket=%d",
			option.Username,
			option.Password,
			option.Addr,
			key,
			0,
		))
		if err != nil {
			panic(fmt.Errorf("db '%s' error: %v", key, err))
		}
		//是否开启日志
		DBMap[key].LogMode(debug)
		//tracemid.GormV1TraceInitialize(DBMap[key])
	}
	return DBMap
}

// CloseDB 关闭数据库
func CloseDB(option interface{}) {
	switch option.(type) {
	case Option:
		_closeDB(option.(Option))
	case []Option:
		for _, dbConfig := range option.([]Option) {
			_closeDB(dbConfig)
		}
	}
}

func _closeDB(option Option) {
	for _, name := range option.DBNames {
		key := option.DatabasePrefix + name
		err := DBMap[key].Close()
		if err != nil {
			log.Fatalf("close err: %v", err)
		}
	}
}
