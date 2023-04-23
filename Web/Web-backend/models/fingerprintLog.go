package models

import (
	"gorm.io/gorm"
)

type FingerprintLog struct {
	gorm.Model
	FingerId string `gorm:"column:finger_id;type:varchar(255);" json:"finger_id"` //指纹id
	WorkNum  string `gorm:"column:work_num;type:varchar(255);" json:"work_num"`   //指纹用户id
	Status   string `gorm:"column:status;type:varchar(255);" json:"status"`       //指纹认证状态
}

func (table *FingerprintLog) TableName() string {
	return "fingerprintLog"
}

func GetFingerprintLogList() *gorm.DB {
	return DB.Find(new(FingerprintLog))

}
