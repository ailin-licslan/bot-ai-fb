package logic

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// 错误提示 & http client
var (
	client                  = fasthttp.Client{}
	errUnknownWebHookObject = errors.New("unknown web hook object")
	errNoMessageEntry       = errors.New("there is no message entry")
	errNoXSignHeader        = errors.New("there is no x-sign header")
	errInvalidXSignHeader   = errors.New("invalid x-sign header")
)

// 常量
const (

	/**
	verifyToken :
	A string that we grab from the Verify Token field in your app's App Dashboard.
	You will set this string when you complete the Webhooks configuration settings steps
	*/
	verifyToken         = "xxx_token"                               //填写真正fb dashboard 填写的保持一致
	appSecret           = "app_secret"                              //填写真正定义的保持一致
	accessToken         = "pls fill in the token from fb dashboard" //与 fb dashboard 的保持一致
	headerNameXSign     = "X-Hub-Signature-256"
	signaturePrefix     = "sha256="
	messageTypeResponse = "RESPONSE"
	// https://developers.facebook.com/docs/messenger-platform/send-messages/#messaging_types
	apiPath               = "https://graph.facebook.com/v12.0/me/messages"
	defaultRequestTimeout = 10 * time.Second
)

// DoTask 请求入口 公开方法 首字母大写
func DoTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		//GET request
		verifyGet(w, r)
		return
	} else {
		//POST request
		reqTaskPost(w, r)
		return
	}

}

// reqTaskPost 处理请求的逻辑
func reqTaskPost(w http.ResponseWriter, r *http.Request) {

	//授权验证
	err := authorize(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		log.Println("authorize", err)
		return
	}

	//获取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		log.Println("read webhook body", err)
		return
	}

	//讲请求数据转化对应的结构体对象
	wr := WebHookRequest{}
	err = json.Unmarshal(body, &wr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		log.Println("unmarshal request", err)
		return
	}

	//处理消息逻辑
	err = mesTaskWebHookExecReq(wr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal"))
		log.Println("handle webhook request", err)
		return
	}

	// Facebook waits for the constant message to get that everything is OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("EVENT_RECEIVED"))
}

// verifyGet 授权校验
func verifyGet(w http.ResponseWriter, r *http.Request) {

	//GET https://www.your-clever-domain-name.com/webhooks?
	//  hub.mode=subscribe&
	//  hub.verify_token=mytoken&
	//  hub.challenge=1158201444

	/**
	Whenever your endpoint receives a verification request, it must:
	Verify that the hub.verify_token value matches the string you set in the
	Verify Token field when you configure the Webhooks product in your App Dashboard
	(you haven't set up this token string yet). Respond with the hub.challenge value.
	*/

	if verifyToken != r.URL.Query().Get("hub.verify_token") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error, wrong validation token"))
		return
	}

	//验证成功
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.Query().Get("hub.challenge")))
}

// authorize 授权校验
func authorize(r *http.Request) error {

	signature := r.Header.Get(headerNameXSign)

	if !strings.HasPrefix(signature, signaturePrefix) {
		return errNoXSignHeader
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	// We read the request body and now it's empty. We have to rewrite it for further reads.
	err = r.Body.Close()
	if err != nil {
		return err
	} //nolint:err check
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	validSignature, err := signatureIsOk(signature, body)
	if err != nil {
		return fmt.Errorf("is valid signature: %w", err)
	}
	if !validSignature {
		return errInvalidXSignHeader
	}

	return nil
}

/**
我们使用 SHA256 签名来签署所有事件通知负载，并在请求的 X-Hub-Signature-256 标头中添加签名，而且会在前面加上
sha256=。验证负载并非强制要求，但我们建议您这样做。
如要验证负载，请执行以下操作：
使用负载和您应用的应用密钥生成 SHA256 签名。
将您的签名与 X-Hub-Signature-256 标头中的签名（sha256= 之后的所有内容）进行对比。如果签名一致，则表示负载真实有效。
请求示例
POST / HTTPS/1.1 Host: your-clever-domain-name.com/webhooks Content-Type: application/json X-Hub-Signature-256: sha256={super-long-SHA256-signature} Content-Length: 311 { "entry": [ { "time": 1520383571, "changes": [ { "field": "photos", "value": { "verb": "update", "object_id": "10211885744794461" } } ], "id": "10210299214172187", "uid": "10210299214172187" } ], "object": "user" }
*/

// signatureIsOk 验证签名
// signatureIsOk https://developers.facebook.com/docs/graph-api/webhooks/getting-started/#validate-payloads
func signatureIsOk(signature string, body []byte) (bool, error) {

	//解析实际签名
	actualSign, err := hex.DecodeString(signature[len(signaturePrefix):])

	if err != nil {
		return false, fmt.Errorf("decode string: %w", err)
	}

	//期望的签名
	expectedSign := signBody(body)

	//判断是否相等
	return hmac.Equal(expectedSign, actualSign), nil
}

// singBody 签名
func signBody(body []byte) []byte {
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Reset()
	h.Write(body)
	return h.Sum(nil)
}

// SignHeadMakingTest 签名生成测试
func SignHeadMakingTest() string {

	//密钥
	appSecretKey := "It is a Secret to Everybody"

	//请求body
	payload := "AI bot"

	h := hmac.New(sha256.New, []byte(appSecretKey))

	h.Write([]byte(payload))

	return signaturePrefix + hex.EncodeToString(h.Sum(nil))
}

// mesTaskWebHookExecReq 处理消息体参数
func mesTaskWebHookExecReq(r WebHookRequest) error {

	if r.Object != "page" {
		return errUnknownWebHookObject
	}

	for _, we := range r.Entry {
		err := mesTaskWebHookExec(we)
		if err != nil {
			return fmt.Errorf("handle webhook request entry: %w", err)
		}
	}

	return nil
}

// mesTaskWebHookExec 处理消息体
func mesTaskWebHookExec(we WebHookRequestEntry) error {

	// Facebook claims that the arr always contains a single item but we don't trust them :)
	if len(we.Messaging) == 0 {
		return errNoMessageEntry
	}

	// message action
	em := we.Messaging[0]
	if em.Message != nil {
		err := mesTaskExec(em.Sender.ID, em.Message.Text)
		if err != nil {
			return fmt.Errorf("handle message: %w", err)
		}
	}

	return nil
}

// mesTaskExec 这里我们需要使用到相关的NLP服务或者自己有能力处理也行 根据用户的问题给予相当智能化聪明的回答 (chatGPT/AutoGPT...)
func mesTaskExec(recipientID, msgText string) error {

	/**
	@TODO It's not my bread and butter in AI / NLP stuff implement by my own, need to invoke AI platform get data
	@TODO It's a simple & dumb bot for facebook messenger right now,
	@TODO to make it smarter and have more interactions with the user.
	@TODO We need to use a NLP backend server like
	@TODO https://openai.com/ (chatGPT) or,
	@TODO https://wit.ai/ (Facebook) or,
	@TODO https://www.motion.ai/ etc.
	@TODO And that will be the subject we need to explore
	*/

	msgText = strings.TrimSpace(msgText)

	var responseText string

	switch msgText {

	// Base on our custom cases to handle Specific Scenarios
	case "go":
		responseText = "I'm going to use goland + NPL like chatGPT to implement AI BOT"

	case "AI":
		responseText = "Let's invoke chatGPT function to give a smart response to users"

	//@FIXME AI ability NLP SERVER / chatGPT call service with those platform
	default:
		responseText = "How can i help you? the response come from NLP server(chatGPT/wit.ai...), " +
			"invoke AI platform interface and get result then invoke sendReqToFaceBookApi back data to user"

	}
	return getResponds(context.TODO(), recipientID, responseText)

}

// getResponds 请求结果获取
func getResponds(ctx context.Context, recipientID, msgText string) error {
	return sendReqToFaceBookApi(ctx, apiPath, setSendReq(recipientID, msgText))
}

// setSendReq 组装请求参数
func setSendReq(recipientID string, msgText string) SendMessageRequest {
	return SendMessageRequest{
		MessagingType: messageTypeResponse,
		RecipientID: MessageRecipient{
			ID: recipientID,
		},
		Message: Message{
			Text: msgText,
		},
	}
}

// sendReqToFaceBookApi 调用 fb api
func sendReqToFaceBookApi(ctx context.Context, reqURI string, reqBody interface{}) error {

	//调试测试使用 使用真正的在fb申请的 access token 这里

	if accessToken == "zzz_accessToken" {
		return nil
	}

	//组装参数
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(fmt.Sprintf("%s?access_token=%s", reqURI, accessToken))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Add("Content-Type", "application/json")
	body, err := json.Marshal(&reqBody)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	req.SetBody(body)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	dl, ok := ctx.Deadline()
	if !ok {
		dl = time.Now().Add(defaultRequestTimeout)
	}
	err = client.DoDeadline(req, res, dl)
	if err != nil {
		return fmt.Errorf("do deadline: %w", err)
	}

	//返回结果
	resp := APIResponse{}
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("response error: %s", resp.Error.Error())
	}
	if res.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("unexpected rsponse status %d", res.StatusCode())
	}

	return nil
}

// ACCESS TOKEN 获取方式
//You can also use the Graph API Explorer  to send the request to subscribe your Page to a Webhooks field.
//1.Select your app in the Application dropdown menu.
//2.Click the Get Token dropdown and select Get User Access Token, then choose the pages_manage_metadata permission. This will exchange your app token for a User access token with the pages_manage_metadata permission granted.
//3.Click Get Token again and select your Page. This will exchange your User access token for a Page access token.
//4.Change the operation method by clicking the GET dropdown menu and selecting POST.
//5.Replace the default me?fields=id,name query with the Page's id followed by /subscribed_apps, then submit the query.

// VERIFY TOKEN 获取方式
//Validating Verification Requests
//Whenever your endpoint receives a verification request, it must:
//Verify that the hub.verify_token value matches the string you set in the Verify Token field when you configure the Webhooks product in your App Dashboard (you haven't set up this token string yet).
//Respond with the hub.challenge value.

//1.Create a Facebook app: https://developers.facebook.com/apps.
//This requires a Facebook Developer account and a Facebook page.

//2. Inside the app, find “Products” on the left menu and click the + icon.
//Find Messenger and click “Set up”. Find “Webhooks” on the page and click “Add callback URL”.

//3. Set the Callback URL to the Ngrok address we set earlier. Input “secret_token” in the second field,
//as this was hardcoded in our basic server setup. If you changed it, set it to what you changed it for.

//4. Connect it to your page by clicking the “Add Pages” button. If you don’t have a page, create one.
//When it’s set up, click the “Generate Token” button next to your page, and save the token it gives
//you for later. To directly work with our code example, export it like so:
//export FACEBOOK_ACCESS_TOKEN="<INPUT YOUR TOKEN HERE>"

//5. Configure page webhook by scrolling back down to “webhooks”. Select the “messages” for now.

//6. Enable Natural language processing from Facebook. We will use this for a really simple matter later.

//7. Test your integration by messaging YOUR page on https://messenger.com (write literally anything). Check your webserver, and you should have a return as seen above.
