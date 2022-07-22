package db_models

import (
	"github.com/jinzhu/gorm"
	"go-ginApp/src/main/pkg/common/trace/trace_service"
	"go-ginApp/src/main/pkg/utils/dbtool"
	"go-ginApp/src/main/pkg/utils/httptool"
	"time"
)

type BookDB struct {
	dbtool.Extra
	dbtool.Model

	DBId          int32     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DBBookName    string    `gorm:"column:book_name" json:"book_name"`
	DBAuthor      string    `gorm:"column:author" json:"author"`
	DBPrice       int32     `gorm:"column:price" json:"price"`
	DBDescribe    string    `gorm:"column:describe" json:"describe"`
	DBReleaseDate time.Time `gorm:"column:release_date" json:"release_date"`
	DBStatus      int32     `gorm:"column:status" json:"status"`
}

func (BookDB) TableName() string {
	return "demo.book"
}

func (d *BookDB) AfterCreate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlAdd, tx.Model(d).Attrs().Value)
}

func (d *BookDB) AfterUpdate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlUpdate, tx.Model(d).Attrs().Value)
}

func (d *BookDB) AfterDelete(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(d.Extra.Ctx, d.TableName(), httptool.LogTypeMysqlDelete, tx.Model(d).Attrs().Value)
}

func (d *BookDB) Id() *int32 {
	return &d.DBId
}

func (d *BookDB) BookName() *string {
	return &d.DBBookName
}

func (d *BookDB) Author() *string {
	return &d.DBAuthor
}
func (d *BookDB) Describe() *string {
	return &d.DBDescribe
}

func (d *BookDB) Price() *int32 {
	return &d.DBPrice
}

func (d *BookDB) ReleaseDate() *string {
	createAt := d.DBReleaseDate.Format("2006-01-02 15:04:05")
	return &createAt
}

func (d *BookDB) Status() *int32 {
	return &d.DBStatus
}

func (d *BookDB) CreatedAt() *string {
	createAt := d.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &createAt
}

func (d *BookDB) UpdatedAt() *string {
	updateAt := d.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &updateAt
}
