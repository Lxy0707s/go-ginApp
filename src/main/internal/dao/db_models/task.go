package db_models

import (
	"github.com/jinzhu/gorm"
	"go-ginApp/src/main/pkg/common/trace/trace_service"
	"go-ginApp/src/main/pkg/utils/dbtool"
	"go-ginApp/src/main/pkg/utils/httptool"
	"time"
)

type DemoTask struct {
	dbtool.Extra
	dbtool.Model

	DId     int32     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DTitle  string    `gorm:"column:title" json:"title"`
	DAuthor string    `gorm:"column:author" json:"author"`
	DVote   int32     `gorm:"column:vote" json:"vote"`
	DDate   time.Time `gorm:"column:date" json:"date"`
}

func (DemoTask) TableName() string {
	return "demo.task"
}

func (d *DemoTask) AfterCreate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlAdd, tx.Model(d).Attrs().Value)
}

func (d *DemoTask) AfterUpdate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlUpdate, tx.Model(d).Attrs().Value)
}

func (d *DemoTask) AfterDelete(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlDelete, tx.Model(d).Attrs().Value)
}

func (d *DemoTask) Id() *int32 {
	return &d.DId
}

func (d *DemoTask) Title() *string {
	return &d.DTitle
}

func (d *DemoTask) Author() *string {
	return &d.DAuthor
}

func (d *DemoTask) Vote() *int32 {
	return &d.DVote
}

func (d *DemoTask) Date() *string {
	createAt := d.DDate.Format("2006-01-02 15:04:05")
	return &createAt
}

func (d *DemoTask) CreatedAt() *string {
	createAt := d.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &createAt
}

func (d *DemoTask) UpdatedAt() *string {
	updateAt := d.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &updateAt
}
