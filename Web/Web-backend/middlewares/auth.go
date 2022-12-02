package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"mp333player.com/utils"
	"net/http"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("token")
		_, err := utils.AnalyseUserToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK,
				gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "Unauthorized",
				})
			return
		}
		c.Next()
	}
}

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("token")
		_, err := utils.AnalyseAdminToken(auth)
		if err != nil {
			c.Abort()
			log.Printf("%+v", err)
			c.JSON(http.StatusOK,
				gin.H{
					"code": http.StatusUnauthorized,
					"msg":  "Unauthorized",
				})
			return
		}
		c.Next()
	}
}
