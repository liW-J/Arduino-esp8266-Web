package main

import (
	"mp333player.com/router"
)

func main() {
	r := router.Router()
	r.Run("0.0.0.0:8888") // 监听并在 0.0.0.0:8888 上启动服务
	//data := make(map[string]interface{})
	//data["openid"] = "123458"
	//data["password"] = "QTD8LmeQ2Bue0oKnyNDK"
	//bytesData, _ := json.Marshal(data)
	//resp, _ := http.Post("http://wx.sends.cc/d46c6Zd6wyEto4gWqyfx/getstuid", "application/json", bytes.NewReader(bytesData))
	//defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
}
