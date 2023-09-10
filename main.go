package main

import (
	"bot-ai-fb/logic"
	"fmt"
	"log"
	"net/http"
	"os"
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

	http.HandleFunc("/", logic.Test)

	//SERVER LOGIC

	//0.GET  fb req verify
	//1.POST send message  [user--> fb --> go app --> NPL back server]
	//2.then response to user [NPL back server --> go app--> fb --> user]
	//3.With NLP SERVER (NPL BACK SERVER) part not implement yet, We can refer to Link https://zhuanlan.zhihu.com/p/611290902 invoke chatGPT API
	//4.If we want test it on local, we could also use a tool like ngrok.
	//It basically creates a secure tunnel on your local machine along with
	//a public URL you can use for browsing your local server. Keep in mind,
	//to use your bot in production, you need to use a real IaaS like AWS, Ali-cloud, QiNiu-cloud, etc

	http.HandleFunc("/webhook", logic.DoTask)

	//GET PORT FROM ENV
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
