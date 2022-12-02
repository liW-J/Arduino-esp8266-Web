package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	WorkNum  string `json:"workNum"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}
type AdminClaims struct {
	Name     string `json:"name"`
	PassWord string `json:"password"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var userKey = []byte("wqrggyyds")
var adminKey = []byte("pjyggyyds")

// GenerateToken
// 生成 token
func GenerateUserToken(workNum string, userName string) string {
	UserClaim := &UserClaims{
		WorkNum:        workNum,
		UserName:       userName,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(userKey)
	if err != nil {
		return ""
	}
	return tokenString
}

func GenerateAdminToken(name string, password string) string {
	AdminClaim := &AdminClaims{
		Name:           name,
		PassWord:       password,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AdminClaim)
	tokenString, err := token.SignedString(adminKey)
	if err != nil {
		return ""
	}
	return tokenString
}

// AnalyseToken
// 解析 token
func AnalyseUserToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return userKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

func AnalyseAdminToken(tokenString string) (*AdminClaims, error) {
	adminClaim := new(AdminClaims)
	claims, err := jwt.ParseWithClaims(tokenString, adminClaim, func(token *jwt.Token) (interface{}, error) {
		return adminKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return adminClaim, nil
}
