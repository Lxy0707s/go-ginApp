package db_models

import (
	"github.com/jinzhu/gorm"
	"go-ginApp/src/main/pkg/common/trace/trace_service"
	"go-ginApp/src/main/pkg/utils/dbtool"
	"go-ginApp/src/main/pkg/utils/httptool"
)

type UserDB struct {
	dbtool.Extra
	dbtool.Model

	DBId       int32  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	DBUserName string `gorm:"column:user_name" json:"user_name"`
	DBToken    string `gorm:"column:token" json:"token"`
	DBEmail    string `gorm:"column:email" json:"email"`
	DBPassword string `gorm:"column:password" json:"password"`
}

func (UserDB) TableName() string {
	return "demo.user"
}

func (u *UserDB) AfterCreate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(u.Extra.Ctx, u.TableName(), httptool.LogTypeMysqlAdd, tx.Model(u).Attrs().Value)
}

func (u *UserDB) AfterUpdate(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(u.Extra.Ctx, u.TableName(), httptool.LogTypeMysqlUpdate, tx.Model(u).Attrs().Value)
}

func (u *UserDB) AfterDelete(tx *gorm.DB) (err error) {
	return trace_service.ApiLogTrace(u.Extra.Ctx, u.TableName(), httptool.LogTypeMysqlDelete, tx.Model(u).Attrs().Value)
}

func (u *UserDB) Id() *int32 {
	return &u.DBId
}

func (u *UserDB) UserName() *string {
	return &u.DBUserName
}

func (u *UserDB) Email() *string {
	return &u.DBEmail
}
func (u *UserDB) Password() *string {
	return &u.DBPassword
}

func (u *UserDB) Token() *string {
	return &u.DBToken
}

func (u *UserDB) CreatedAt() *string {
	createAt := u.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &createAt
}

func (u *UserDB) UpdatedAt() *string {
	updateAt := u.Model.CreatedAt.Format("2006-01-02 15:04:05")
	return &updateAt
}
