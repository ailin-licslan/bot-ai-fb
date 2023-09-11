package midleware

import "github.com/gin-gonic/gin"

func X(c *gin.Context) {

	println("xxx this is 中间件测试")
}
