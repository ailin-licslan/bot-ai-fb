package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var adminList = gin.H{
	"LinS": gin.H{"email": "lin1@sina.com", "phone": "13125172881"},
	"LinQ": gin.H{"email": "lin2@sina.com", "phone": "13125172882"},
	"LinK": gin.H{"email": "lin3@sina.com", "phone": "13125172883"},
	"LinM": gin.H{"email": "lin4@sina.com", "phone": "13125172884"},
}

func CheckAuth(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	if res, ok := adminList[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "res": res})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "res": "no data"})
	}
}

func NameList() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		//这个可以从数据库里面取出来  用户名+密码
		"LinA": "123",
		"LinQ": "456",
		"LinK": "789",
	})
}
