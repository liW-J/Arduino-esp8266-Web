package main

import (
	"mp333player.com/router"
)

func main() {
	r := router.Router()
	r.Run("0.0.0.0:8888") // 监听并在 0.0.0.0:8888 上启动服务
}
