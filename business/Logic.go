package business

import (
	common "bot-ai-fb/common"
	Db "bot-ai-fb/database"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// LoginTask  return user err
func LoginTask(u *common.LoginData) (user *common.User, err error) {

	// 记录一下原始密码(用户登录的密码)
	originPassword := u.Password

	sqlStr := "SELECT id, user_name, password FROM USER WHERE user_name = ?"

	_, err = Db.Dbs.Query(sqlStr, u.UserName)
	// 查询数据库出错
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// 用户不存在
	if err == sql.ErrNoRows {
		return nil, errors.New(common.ErrorUserNotExit)
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return nil, errors.New(common.ErrorPasswordWrong)
	}

	accessToken, refreshToken, err := common.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return user, err
}

const secret = "LICSLAN"

// encryptPassword 对密码进行加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

func SignUp(p *common.RegisterForm) (error error) {
	// 1、判断用户存不存在
	err := CheckUserExist(p.UserName)
	if err != nil {
		// 数据库查询出错
		return err
	}

	// 2、生成UID
	userId, err := common.GetID()
	if err != nil {
		return common.ErrorGenIDFailed
	}
	// 构造一个User实例
	u := common.User{
		UserID:   userId,
		UserName: p.UserName,
		Password: p.Password,
		Email:    p.Email,
		Gender:   p.Gender,
	}
	// 3、保存进数据库
	return SaveUser(u)
}

func CheckUserExist(username string) (error error) {

	sqlStr := `SELECT count(user_id) FROM user WHERE username = ?`

	query, err := Db.Dbs.Query(sqlStr, username)
	if err != nil {
		return err
	}

	if query != nil {
		return errors.New(common.ErrorUserExit)
	}
	return
}

// SaveUser 注册业务-向数据库中插入一条新的用户
func SaveUser(user common.User) (error error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	// 执行SQL语句入库
	sqlStr := `INSERT INTO user(user_id,username,password,email,gender) VALUES(?,?,?,?,?)`
	_, err := Db.Dbs.Exec(sqlStr, user.UserID, user.UserName, user.Password, user.Email, user.Gender)
	return err
}
