package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB = Init()

func Init() *gorm.DB {
	dsn := "name:password@tcp(your_ip:3306)/Arduino_esp8266_Web?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm.Init ERROR :", err)
	}

	//自动创建表
	err = db.AutoMigrate(&User{}, &FingerprintLog{}, &TemperatureLog{}, &OperateLog{})
	if err != nil {

		log.Println("AutoMigrate ERROR :", err)
	}

	//创建记录
	u1 := User{WorkNum: "123456", UserName: "小王", Finger1: "1"}
	t1 := TemperatureLog{TemperatureData: "36.78", WorkNum: "123456", UserName: "小王"}
	f1 := FingerprintLog{WorkNum: "123456", Status: "认证成功", FingerId: "1"}
	//s1 := Song{SenderName: "123", SenderStuNum: "2026311013", ReceiverName: "123", SongId: "123", PhoneNum: "123", SchoolDistrict: "123", SearchPath: "123", BlessingWords: "123", BroadcastDate: time.Now().UTC()}
	db.Create(&u1)
	db.Create(&t1)
	db.Create(&f1)

	return db

}
