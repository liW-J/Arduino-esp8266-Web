package models

import (
	"gorm.io/gorm"
)

type TemperatureLog struct {
	gorm.Model
	TemperatureData string `gorm:"column:temperature_data;type:varchar(255);" json:"fingerprint_data"` //温度数据
	UserName        string `gorm:"column:user_name;type:varchar(255);" json:"user_name"`               //指纹用户姓名
	WorkNum         string `gorm:"column:work_num;type:varchar(255);" json:"work_num"`                 //工号
}

func (table *TemperatureLog) TableName() string {
	return "temperatureLog"
}

func GetTemperatureLogList() *gorm.DB {
	return DB.Find(new(TemperatureLog))

}
