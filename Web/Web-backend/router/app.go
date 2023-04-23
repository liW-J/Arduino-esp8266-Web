package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "mp333player.com/docs"
	"mp333player.com/middlewares"
	"mp333player.com/service"
	"mp333player.com/service/ws"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()
	//配置swagger
	r.Use(Cors())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//路由规则

	//用户

	hub := ws.NewHub()
	//开启监听通道
	go hub.Run()
	//管理员
	authAdmin := r.Group("/admin")
	authAdmin.POST("/login", service.AdminLogin)
	authAdmin.GET("/fingerLog", middlewares.AuthAdminCheck(), service.GetFingerList)
	authAdmin.GET("/tempLog", middlewares.AuthAdminCheck(), service.GetTempList)
	authAdmin.GET("/userList", middlewares.AuthAdminCheck(), service.GetUserList)
	authAdmin.GET("/operateLog", middlewares.AuthAdminCheck(), service.GetOperateList)
	authAdmin.POST("/userTempLog", middlewares.AuthAdminCheck(), service.GetUserTempLog)
	authAdmin.POST("deleteUser", middlewares.AuthAdminCheck(), service.DeleteUser)
	//用户注册
	authAdmin.POST("/signUp", middlewares.AuthAdminCheck(), func(c *gin.Context) {
		service.SignUp(c, hub)
	})

	//用户
	authUser := r.Group("/user")
	//建立ws
	authUser.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	//验证用户指纹
	authUser.GET("/login", func(c *gin.Context) {
		service.UserLoginByFinger(c, hub)
	})
	//记录体温
	authUser.POST("/recordTemp", middlewares.AuthUserCheck(), func(c *gin.Context) {
		service.RecordTemp(c, hub)
	})
	//新增指纹
	authUser.POST("/addFinger", middlewares.AuthUserCheck(), func(c *gin.Context) {
		service.AddFinger(c, hub)
	})
	//更新指纹
	authUser.POST("/updateFinger", middlewares.AuthUserCheck(), func(c *gin.Context) {
		service.UpdateFinger(c, hub)
	})
	return r
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}
