package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"mp333player.com/models"
	"mp333player.com/service/ws"
	"mp333player.com/utils"
	"mp333player.com/utils/e"
	"net/http"
	"strconv"
	"time"
)

// UserLoginByFinger
// @Tags 用户
// @Summary 用户通过指纹身份
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/login [get]
func UserLoginByFinger(c *gin.Context, h *ws.Hub) {

	//连接ws相关
	sendMsg := ws.ReplyMsg{
		Code: e.SUCCESS,
		Data: "Verify fingerprint......",
	}
	message, _ := json.Marshal(sendMsg)
	fmt.Println(len(h.Clients))
	for c := range h.Clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	//ws发送完成
	//15s延时，保证指纹已成功采集
	time.Sleep(time.Duration(15) * time.Second)

	//获取最新指纹记录
	fingerData := make([]*models.FingerprintLog, 0)
	fingerTx := models.GetFingerprintLogList()
	_ = fingerTx.Last(&fingerData).Error
	//检验是不是有存入新数据
	now := time.Now()
	ss, _ := time.ParseDuration("-15s")
	ss15 := now.Add(ss)
	if ss15.After(fingerData[0].CreatedAt) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ws连接创建失败!请重试！",
			"data": "",
		})
		return
	}

	finger := fingerData[0].FingerId
	data := make([]*models.User, 0)
	tx := models.GetUserList()
	//判断是否存在
	var cnt int64
	err := tx.Where("finger1 = ? ", finger).Or("finger2 = ?", finger).Count(&cnt).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "无效的指纹",
			"data": "",
		})
		return
	}
	if fingerData[0].FingerId == "0" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证超时！请及时将指纹置于传感器上！",
			"data": "",
		})
		return
	}
	if cnt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户不存在，请联系管理员录入指纹！",
			"data": "",
		})
		return
	}
	fingerData[0].WorkNum = data[0].WorkNum
	fingerData[0].Status = "认证成功"
	err = fingerTx.Save(&fingerData).Error
	if err != nil {
		log.Printf("save err %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统故障，请重试！",
			"data": "",
		})
		return
	}
	//验证通过生成token
	token := utils.GenerateUserToken(data[0].WorkNum, data[0].UserName)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "认证成功",
		"data": token,
	})

}

// RecordTemp
// @Tags 用户
// @Summary 记录用户体温
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/recordTemp [post]
func RecordTemp(c *gin.Context, h *ws.Hub) {
	//验证用户信息
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseUserToken(auth)
	//建立ws连接
	sendMsg := ws.ReplyMsg{
		Code: e.SUCCESS,
		Data: "Record Temperature......",
	}
	message, _ := json.Marshal(sendMsg)
	fmt.Println(len(h.Clients))
	for c := range h.Clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	//ws发送完成
	//5s延时，保证指纹已成功采集
	time.Sleep(time.Duration(5) * time.Second)

	//获取最新体温记录
	tempData := make([]*models.TemperatureLog, 0)
	tempTx := models.GetTemperatureLogList()
	_ = tempTx.Last(&tempData).Error
	//检验是不是有存入新数据
	now := time.Now()
	ss, _ := time.ParseDuration("-5s")
	ss15 := now.Add(ss)
	if ss15.After(tempData[0].CreatedAt) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ws连接创建失败!请重试！",
			"data": "",
		})
		return
	}
	tempData[0].UserName = resp.UserName
	tempData[0].WorkNum = resp.WorkNum
	err := tempTx.Save(&tempData).Error
	if err != nil {
		log.Printf("save err %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统故障，请重试！",
			"data": "",
		})
		return
	}
	v, _ := strconv.ParseFloat(tempData[0].TemperatureData, 32)
	//体温小于37.3为正常
	if v <= 37.3 {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "测温成功，体温正常",
			"data": gin.H{
				"temperature": tempData[0].TemperatureData,
			},
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "体温异常！请联系工作人员！",
			"data": "",
		})
		return
	}
}

// AddFinger
// @Tags 用户
// @Summary 添加指纹
// @Param token header string true "token"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/addFinger [post]
func AddFinger(c *gin.Context, h *ws.Hub) {
	//验证用户信息
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseUserToken(auth)
	data := make([]*models.User, 0)
	tx := models.GetUserList()
	err := tx.Where("work_num = ?", resp.WorkNum).Find(&data).Error
	if err != nil {
		log.Printf("find user list err: %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户信息获取失败！请重试！",
			"data": "",
		})
		return
	}
	if data[0].Finger2 != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "已存在指纹2",
			"data": "",
		})
		return
	}

	//建立ws连接
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
	//5s延时，保证指纹已成功采集
	time.Sleep(time.Duration(15) * time.Second)

	//获取最新体温记录
	fingerData := make([]*models.FingerprintLog, 0)
	tempTx := models.GetFingerprintLogList()
	_ = tempTx.Last(&fingerData).Error
	//检验是不是有存入新数据
	now := time.Now()
	ss, _ := time.ParseDuration("-15s")
	ss15 := now.Add(ss)
	log.Println(ss15)
	log.Println(fingerData[0].CreatedAt)
	if ss15.After(fingerData[0].CreatedAt) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ws连接创建失败!请重试！",
			"data": "",
		})
		return
	}
	if fingerData[0].FingerId == "0" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "添加指纹超时！请及时将指纹置于传感器上！",
			"data": "",
		})
	}
	data[0].Finger2 = fingerData[0].FingerId
	err = tx.Save(&data).Error
	if err != nil {
		log.Printf("save err %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "无法添加指纹信息！",
			"data": "",
		})
		return
	}
	newOperate := models.OperateLog{
		Role:    resp.WorkNum,
		Operate: "添加指纹",
		Object:  "指纹2",
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
		"msg":  "成功添加指纹2",
		"data": "",
	})
	return

}

// UpdateFinger
// @Tags 用户
// @Summary 更新指纹
// @Param token header string true "token"
// @Param updateInfo body models.UpdateFingerInfo true "updateInfo"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/updateFinger [post]
func UpdateFinger(c *gin.Context, h *ws.Hub) {
	//验证用户信息
	auth := c.GetHeader("token")
	resp, _ := utils.AnalyseUserToken(auth)
	data := make([]*models.User, 0)
	tx := models.GetUserList()
	err := tx.Where("work_num = ?", resp.WorkNum).Find(&data).Error
	if err != nil {
		log.Printf("find user list err: %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户信息获取失败！请重试！",
			"data": "",
		})
		return
	}

	//建立ws连接
	sendMsg := ws.ReplyMsg{
		Code: e.SUCCESS,
		Data: "Update fingerprint......",
	}
	message, _ := json.Marshal(sendMsg)
	fmt.Println(len(h.Clients))
	for c := range h.Clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	//ws发送完成
	//5s延时，保证指纹已成功采集
	time.Sleep(time.Duration(15) * time.Second)

	//获取最新体温记录
	fingerData := make([]*models.FingerprintLog, 0)
	fingerTx := models.GetFingerprintLogList()
	_ = fingerTx.Last(&fingerData).Error
	//检验是不是有存入新数据
	now := time.Now()
	ss, _ := time.ParseDuration("-15s")
	ss15 := now.Add(ss)
	if ss15.After(fingerData[0].CreatedAt) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ws连接创建失败!请重试！",
			"data": "",
		})
		return
	}
	if fingerData[0].FingerId == "0" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "录入指纹超时！请及时将指纹置于传感器上！",
			"data": "",
		})
	}
	//更新指纹
	json := models.UpdateFingerInfo{}
	c.BindJSON(&json)
	if json.UpdateNum == 1 {
		data[0].Finger1 = fingerData[0].FingerId
	} else if json.UpdateNum == 2 {
		data[0].Finger2 = fingerData[0].FingerId
	}

	err = tx.Save(&data).Error
	if err != nil {
		log.Printf("save err %+v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "无法更新指纹信息！",
			"data": "",
		})
		return
	}
	newOperate := models.OperateLog{
		Role:    resp.WorkNum,
		Operate: "更新指纹",
		Object:  "指纹" + string(json.UpdateNum),
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
		"msg":  "指纹更新成功",
		"data": "",
	})
	return
}
