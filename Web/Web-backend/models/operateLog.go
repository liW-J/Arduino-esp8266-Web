package models

import (
	"gorm.io/gorm"
)

type OperateLog struct {
	gorm.Model
	Role    string `gorm:"column:role;type:varchar(255);" json:"role"`       //操作者
	Operate string `gorm:"column:Operate;type:varchar(255);" json:"operate"` // 操作事件
	Object  string `gorm:"column:Object;type:varchar(255);" json:"object"`   // 操作对象
}

func (table *OperateLog) TableName() string {
	return "operateLog"
}

func GetOperateLogList() *gorm.DB {
	return DB.Find(new(OperateLog))
}
