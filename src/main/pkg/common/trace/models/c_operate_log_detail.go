package models

import (
	"go-ginApp/src/main/pkg/utils/datetool"
	"time"
)

type OperateLogDetail struct {
	Id            int32     `json:"id" gorm:"primary_key"`
	RequestId     string    `json:"request_id"`
	ApiId         int32     `json:"api_id"`
	OpTableName   string    `json:"op_table_name"`
	TableAction   int32     `json:"table_action"`
	ChangeContent string    `json:"change_content"`
	CreatedAt     time.Time `json:"created_at"`
}

func (OperateLogDetail) TableName() string {
	// 此为公用表，不指定具体数据库名
	return "c_operate_log_detail"
}

type OperateLogDetailGType struct {
	*OperateLogDetail
}

func (g *OperateLogDetailGType) Id() *int32 {
	return &g.OperateLogDetail.Id
}

func (g *OperateLogDetailGType) RequestId() *string {
	return &g.OperateLogDetail.RequestId
}

func (g *OperateLogDetailGType) ApiId() *int32 {
	return &g.OperateLogDetail.ApiId
}

func (g *OperateLogDetailGType) OpTableName() *string {
	return &g.OperateLogDetail.OpTableName
}

func (g *OperateLogDetailGType) TableAction() *int32 {
	return &g.OperateLogDetail.TableAction
}

func (g *OperateLogDetailGType) ChangeContent() *string {
	return &g.OperateLogDetail.ChangeContent
}

func (g *OperateLogDetailGType) CreatedAt() *string {
	CreatedAt := datetool.TimestampToTime(g.OperateLogDetail.CreatedAt.Unix())
	return &CreatedAt
}
