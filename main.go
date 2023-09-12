package main

import (
	"bot-ai-fb/database"
	"bot-ai-fb/ginweb"
	"bot-ai-fb/logic"
	"bot-ai-fb/midleware"
	"bot-ai-fb/redis"
	settings "bot-ai-fb/setting"
	"github.com/gin-gonic/gin"
)

/**
参考官方文档
//So a user sends a message to the Facebook Messenger bot. On message,
//Facebook sends webhook to our server (Golang App). The server handles
//the message and responds to the user by Facebook Messenger API.
//https://developers.facebook.com/docs/messenger-platform/webhooks#requirements
//https://developers.facebook.com/docs/messenger-platform/webhooks#configure-webhooks-product
*/

func main() {

	//Gin web 方式启动   gin.Default() 默认内部调用New Logger Recovery中间件
	router := gin.Default()
	//中间件的使用
	router.Use(midleware.X)
	//加载配置文件
	settings.Init()
	//数据库连接测试  这个放在前面哈 !!!  这里连接配置可以放到配置文件里面或者以后有配置中心也行
	//database.ConnectDb()
	configMySQL := settings.Conf.MySQLConfig
	database.InitDB(configMySQL)
	//程序退出 关闭数据库连接
	defer database.Close()

	//REDIS 初始化
	configRedis := settings.Conf.RedisConfig
	redis.InitRedis(configRedis)
	//程序退出 关闭redis连接
	defer redis.Close()

	//路由
	router.GET("/", ginweb.Y)
	router.GET("/x", ginweb.Test)

	//RESTFUL   参数解析
	router.GET("/y/:name", ginweb.ParamTest)
	//?name=123 参数解析
	router.GET("/y/", ginweb.ParamTest)

	//INSERT TEST  GET 方便浏览器测试 不用开POSTMAN
	router.GET("/save/:name/:age", ginweb.Save)
	router.GET("/del/:id", ginweb.Delete)

	//REDIS TEST
	router.GET("/redis/:name", ginweb.RedisAdd)

	//面试测试  localhost:9999/webhook  GET/POST
	router.GET("/webhook", logic.DoTaskV2)
	router.POST("/webhook", logic.DoTaskV2)

	//启动服务
	router.Run("127.0.0.1:9999")

	//传统http方式启动  http.ListenAndServe
	//SERVER LOGIC
	//0.GET  fb req verify
	//1.POST send message  [user--> fb --> go app --> NPL back server]
	//2.then response to user [NPL back server --> go app--> fb --> user]
	//3.With NLP SERVER (NPL BACK SERVER) part not implement yet, We can refer to Link https://zhuanlan.zhihu.com/p/611290902 invoke chatGPT API
	//4.If we want test it on local, we could also use a tool like ngrok.
	//It basically creates a secure tunnel on your local machine along with
	//a public URL you can use for browsing your local server. Keep in mind,
	//to use your bot in production, you need to use a real IaaS like AWS, Ali-cloud, QiNiu-cloud, etc

	//http.HandleFunc("/", logic.Test)
	//http.HandleFunc("/webhook", logic.DoTask)
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8888"
	//	log.Printf("Defaulting to port %s", port)
	//}
	//log.Printf("Listening on port %s", port)
	//log.Printf("Open http://localhost:%s in the browser", port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
