package auth

import (
	bus "bot-ai-fb/business"
	"bot-ai-fb/common"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var trans ut.Translator

func Login(c *gin.Context) {

	//1.参数校验
	var uLogin *common.LoginData
	err, done, u := check(c, uLogin)
	if done {
		return
	}

	//2.登录业务
	user, done2 := taskLogin(c, err, u)
	if done2 {
		return
	}

	// 3、返回响应
	getData(c, user)
}

func SignUp(c *gin.Context) {

	// 1.获取请求参数
	var fo *common.SignUpData

	// 2.校验数据有效性
	upCheck, data := signUpCheck(c, fo)
	if upCheck {
		return
	}

	fmt.Printf("fo: %v\n", data)
	// 3.业务处理 —— 注册用户
	if singUpTask(c, data) {
		return
	}

	//返回响应
	common.ResponseSuccess(c, nil)
}

func RefreshToken(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的 Authorization 中，并使用 Bearer 开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		common.ResponseErrorWithMsg(c, common.CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		common.ResponseErrorWithMsg(c, common.CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := common.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}

func singUpTask(c *gin.Context, fo *common.SignUpData) bool {
	if err := bus.SignUp(fo); err != nil {
		zap.L().Error("logic.signup failed", zap.Error(err))
		if err.Error() == common.ErrorUserNotExit {
			common.ResponseError(c, common.CodeUserExist)
			return true
		}
		common.ResponseError(c, common.CodeServerBusy)
		return true
	}
	return false
}

func signUpCheck(c *gin.Context, fo *common.SignUpData) (bool, *common.SignUpData) {

	if err := c.ShouldBindJSON(&fo); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			common.ResponseError(c, common.CodeInvalidParams) // 请求参数错误
			return true, nil
		}
		// validator.ValidationErrors类型错误则进行翻译
		common.ResponseErrorWithMsg(c, common.CodeInvalidParams, common.RemoveTopStruct(errs.Translate(trans)))
		return true, nil // 翻译错误
	}
	return false, fo
}

func getData(c *gin.Context, user *common.User) {
	common.ResponseSuccess(c, gin.H{
		"id":           fmt.Sprintf("%d", user.UserID), //js识别的最大值：id值大于1<<53-1  int64: i<<63-1
		"username":     user.UserName,
		"accessToken":  user.AccessToken,
		"refreshToken": user.RefreshToken,
	})
}

func taskLogin(c *gin.Context, err error, uLogin *common.LoginData) (*common.User, bool) {
	user, err := bus.LoginTask(uLogin)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", uLogin.UserName), zap.Error(err))
		if err.Error() == common.ErrorUserNotExit {
			common.ResponseError(c, common.CodeUserNotExist)
			return nil, true
		}
		common.ResponseError(c, common.CodeInvalidParams)
		return nil, true
	}
	return user, false
}

func check(c *gin.Context, uLogin *common.LoginData) (error, bool, *common.LoginData) {
	err := c.ShouldBindJSON(&uLogin)
	if err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			common.ResponseError(c, common.CodeInvalidParams) // 请求参数错误
			return nil, true, nil
		}
		// validator.ValidationErrors类型错误则进行翻译
		common.ResponseErrorWithMsg(c, common.CodeInvalidParams,
			common.RemoveTopStruct(errors.Translate(trans)))
		return nil, true, nil
	}
	return err, false, uLogin
}
