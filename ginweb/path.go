package ginweb

import (
	"bot-ai-fb/crud"
	"bot-ai-fb/database"
	"bot-ai-fb/model"
	"bot-ai-fb/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": "200",
	})
}

func Y(c *gin.Context) {
	c.JSON(200, gin.H{
		"Hello": "gin-web learn",
		"Go":    "start",
	})
}

func ParamTest(c *gin.Context) {

	//获取参数 xxx/123 get RestFul
	name := c.Param("name")

	//获取参数 xxx?name=123
	name2 := c.Query("name")

	if name != "" {
		c.JSON(200, gin.H{
			"Hello": "gin-web learn",
			"Go":    name,
		})
	} else {
		c.JSON(200, gin.H{
			"Hello": "gin-web learn",
			"Go":    name2,
		})
	}

}

func Save(c *gin.Context) {
	name := c.Param("name")
	age, _ := strconv.Atoi(c.Param("age"))
	user := model.User{
		Name: name, Age: age,
	}
	crud.InsertUser(database.Dbs, user)
	c.JSON(200, user)
}

func Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	crud.DeleteUser(database.Dbs, id)
	c.JSON(200, id)
}

func RedisAdd(c *gin.Context) {
	name := c.Param("name")
	crud.RedisAdd(redis.Con, name)
	c.JSON(200, gin.H{
		"name": name,
	})
}
