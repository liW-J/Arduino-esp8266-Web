package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"log"
	"mp333player.com/config"
	"mp333player.com/models"
	"mp333player.com/service/ws"
	"mp333player.com/utils"
	"mp333player.com/utils/e"
	"net/http"
	"time"
)

// AdminLogin
// @Tags 管理员
// @Summary 管理员登陆
// @Param adminInfo body models.AdminInfo true "管理员登陆"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	json := models.AdminInfo{}
	c.BindJSON(&json)
	name := json.Name
	password := json.Password
	log.Printf(name)
	log.Printf(password)
	if name != config.Name || password != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
			"data": "",
		})
		return
	}
	token := utils.GenerateAdminToken(name, password)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登陆成功",
		"data": gin.H{
			"token": token,
		},
	})

}

// SignUp
// @Tags 管理员
// @Summary 用户注册并录入第一个指纹
// @Param token header string true "token"
// @Param signUpInfo body models.SignUpInfo true "signUpInfo"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/signUp [post]
func SignUp(c *gin.Context, h *ws.Hub) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}

	//连接ws相关
	sendMsg := ws.ReplyMsg{
		Code: e.SUCCESS,
		Data: "Add fingerprint......",
	}
	message, _ := json.Marshal(sendMsg)
	fmt.Println(len(h.Clients))
	for c := range h.Clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	//ws发送完成
	//20s延时，保证指纹已成功采集
	time.Sleep(time.Duration(20) * time.Second)

	//获取最新指纹记录
	fingerData := make([]*models.FingerprintLog, 0)
	fingerTx := models.GetFingerprintLogList()
	_ = fingerTx.Last(&fingerData).Error
	//检验是不是有存入新数据
	now := time.Now()
	ss, _ := time.ParseDuration("-20s")
	ss15 := now.Add(ss)
	if ss15.After(fingerData[0].CreatedAt) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ws连接创建失败!请重试！",
			"data": "",
		})
		return
	}
	json := models.SignUpInfo{}
	c.BindJSON(&json)
	if json.UserName == "" || json.WorkNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空！",
			"data": "",
		})
		return
	}
	user := models.User{
		UserName: json.UserName,
		WorkNum:  json.WorkNum,
		Finger1:  fingerData[0].FingerId,
	}
	err := models.DB.Create(&user).Error
	if err != nil {
		log.Printf("creat User err:%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建用户失败！",
			"data": "",
		})
		return
	}
	newOperate := models.OperateLog{
		Role:    "Admin",
		Operate: "添加新用户",
		Object:  json.WorkNum,
	}
	err = models.DB.Create(&newOperate).Error
	if err != nil {
		log.Printf("creat temperatureLog err:%+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "操作日志保存失败！",
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "创建成功",
		"data": "",
	})

}

// GetUserList
// @Tags 管理员
// @Summary 获取用户列表
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/userList [get]
func GetUserList(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	tx := models.GetUserList()
	data := make([]*models.User, 0)
	err := tx.Find(&data).Error
	if err != nil {
		log.Println("Get user List Error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})

}

// GetFingerList
// @Tags 管理员
// @Summary 获取指纹识别日志
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/fingerLog [get]
func GetFingerList(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	data := make([]*models.FingerprintLog, 0)
	tx := models.GetFingerprintLogList()
	err := tx.Find(&data).Error
	if err != nil {
		log.Println("Get finger List Error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})

}

// GetTempList
// @Tags 管理员
// @Summary 获取温度日志
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/tempLog [get]
func GetTempList(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	data := make([]*models.TemperatureLog, 0)
	tx := models.GetTemperatureLogList()
	err := tx.Find(&data).Error
	if err != nil {
		log.Println("Get temp List Error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})

}

// GetOperateList
// @Tags 管理员
// @Summary 获取操作日志
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/operateLog [get]
func GetOperateList(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	data := make([]*models.OperateLog, 0)
	tx := models.GetOperateLogList()
	err := tx.Find(&data).Error
	if err != nil {
		log.Println("Get operate List Error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})

}

// DeleteUser
// @Tags 管理员
// @Summary 删除用户
// @Param token header string true "token"
// @Param signUpInfo body models.UserInfo true "signUpInfo"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/deleteUser [post]
func DeleteUser(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	json := models.UserInfo{}
	c.BindJSON(&json)
	tx := models.GetUserList()
	data := make([]*models.User, 0)
	err := tx.Where("work_num = ?", json.WorkNum).Find(&data).Error
	if err != nil {
		log.Println("Find user Error", err)
		return
	}
	err = tx.Delete(&data).Error
	if err != nil {
		log.Println("Delete user Error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": "",
	})

}

// GetUserTempLog
// @Tags 管理员
// @Summary 获取用户体温记录
// @Param token header string true "token"
// @Param signUpInfo body models.UserInfo true "signUpInfo"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/userTempLog [post]
func GetUserTempLog(c *gin.Context) {
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseAdminToken(auth)
	if resp.Name != config.Name || resp.PassWord != config.Password {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "你不是管理员",
			"data": "",
		})
		return
	}
	json := models.UserInfo{}
	c.BindJSON(&json)
	tx := models.GetTemperatureLogList()
	data := make([]*models.TemperatureLog, 0)
	err := tx.Where("work_num = ?", json.WorkNum).Find(&data).Error
	if err != nil {
		log.Println("Find user Error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "OK",
		"data": data,
	})

}
