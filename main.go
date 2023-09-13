package main

import (
	Auth "bot-ai-fb/auth"
	"bot-ai-fb/common"
	Db "bot-ai-fb/database"
	Gin "bot-ai-fb/ginweb"
	Log "bot-ai-fb/logger"
	Logic "bot-ai-fb/logic"
	Mid "bot-ai-fb/midleware"
	Redis "bot-ai-fb/redis"
	settings "bot-ai-fb/setting"
	"github.com/gin-gonic/gin"
	Swag "github.com/swaggo/files"
	GinSwag "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
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

	//=====================================加载配置部分========================================

	//Gin web 方式启动   gin.Default() 默认内部调用New Logger Recovery中间件
	router := gin.Default()
	//中间件的使用
	router.Use(
		Mid.X,                            //测试中间件
		Log.GinRecovery(true),            // Recovery 中间件会 recover掉项目可能出现的panic，并使用zap记录相关日志
		Mid.RateLimit(2*time.Second, 40), // 每两秒钟添加十个令牌  全局限流
	)
	//加载配置文件
	settings.Init()
	//数据库连接测试  这个放在前面哈 !!!  这里连接配置可以放到配置文件里面或者以后有配置中心也行
	//database.ConnectDb()
	configMySQL := settings.Conf.MySQLConfig
	Db.InitDB(configMySQL)
	//程序退出 关闭数据库连接
	defer Db.Close()
	//日志初始化
	Log.Init(settings.Conf.LogConfig, settings.Conf.Mode)
	//REDIS 初始化
	configRedis := settings.Conf.RedisConfig
	Redis.InitRedis(configRedis)
	//程序退出 关闭redis连接
	defer Redis.Close()
	// 雪花算法生成分布式ID
	common.Init(1)
	// 初始化翻译器
	common.InitTrans("zh")

	//=====================================测试学习部分========================================

	//静态文件加载
	router.LoadHTMLFiles("templates/index.html") // 加载HTML
	router.Static("/static", "./static")         // 加载静态文件
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// 注册 SWAGGER
	router.GET("/swagger/*any", GinSwag.WrapHandler(Swag.Handler))

	//BasicAuth路由组权限案例   用户授权校验  Auth.NameList()
	groupAuth := router.Group("/admin", Mid.JWTAuth())
	//路由的前缀 base path /v1/x
	groupAuth.GET("/v1/x", Auth.CheckAuth)
	//路由 测试
	groupAuth.GET("/z", Gin.Y)
	groupAuth.GET("/x", Gin.Test)

	//RESTFUL   参数解析  测试
	groupAuth.GET("/y/:name", Gin.ParamTest)
	//?name=123 参数解析  测试
	groupAuth.GET("/y", Gin.ParamTest)

	//INSERT TEST  GET 方便浏览器测试 不用开POSTMAN 测试
	groupAuth.GET("/save/:name/:age", Gin.Save)
	groupAuth.GET("/del/:id", Gin.Delete)

	//REDIS TEST  测试
	router.GET("/redis/:name", Gin.RedisAdd)
	//面试测试  localhost:9999/webhook  GET/POST
	router.GET("/webhook", Logic.DoTaskV2)
	router.POST("/webhook", Logic.DoTaskV2)

	//=====================================业务开发部分========================================

	// 登录注册 token JWT
	router.POST("/Login", Auth.Login)
	router.POST("/SignUp", Auth.SignUp)
	router.GET("/RefreshToken", Auth.RefreshToken) // 刷新accessToken

	// 中间件 下面的接口需要登录后带上JWT去请求的
	//groupAuthToken := router.Group("/admin-api")
	//groupAuthToken.Use(Mid.JWTAuth()) // 应用JWT认证中间件
	//{
	//	//业务开发
	//	groupAuthToken.POST("/post", nil)       // 创建帖子
	//	groupAuthToken.POST("/vote", nil)       // 投票
	//	groupAuthToken.POST("/comment", nil)    // 评论
	//	groupAuthToken.GET("/commentList", nil) // 评论列表
	//	groupAuthToken.GET("/ping", func(c *gin.Context) {
	//		c.String(http.StatusOK, "pong")
	//	})
	//}

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
