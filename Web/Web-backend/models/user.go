package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"column:user_name;type:varchar(255);" json:"user_name"` //用户姓名
	WorkNum  string `gorm:"column:work_num;type:varchar(255);" json:"work_num"`   //工号
	Finger1  string `gorm:"column:finger1;type:varchar(255);" json:"finger1"`     //用户的第一个指纹
	Finger2  string `gorm:"column:finger2;type:varchar(255);" json:"finger2"`     //用户的第二个指纹
}

type UserInfo struct {
	WorkNum string `json:"work_num"` //工号
}

type SignUpInfo struct {
	UserName string `json:"user_name"` //用户姓名
	WorkNum  string `json:"work_num"`  //工号
}

type AdminInfo struct {
	Name     string `json:"name"`     //name
	Password string `json:"password"` //password
}

type AddFingerInfo struct {
	Finger  string `json:"finger"`   //指纹
	WorkNum string `json:"work_num"` //工号
}

type UpdateFingerInfo struct {
	UpdateNum int `json:"update_num"` //更换指纹的序号  1or2
}

func (table *User) TableName() string {
	return "user"
}

func GetUserList() *gorm.DB {
	return DB.Model(new(User))
}
