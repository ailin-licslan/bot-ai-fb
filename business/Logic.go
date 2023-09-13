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

	sqlStr := "SELECT user_id,username,password FROM t_user WHERE username = ?"

	var us []common.User
	rows, err := Db.Dbs.Query(sqlStr, u.UserName)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u common.User
		// sqlx 提供了便捷方法可以将查询结果直接扫描到结构体
		err2 := rows.Scan(&u.UserID, &u.UserName, &u.Password)
		if err2 != nil {
			return nil, err2
		}
		us = append(us, u)
	}

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
	p := us[0].Password
	id := us[0].UserID
	name := us[0].UserName
	if p != password {
		return nil, errors.New(common.ErrorPasswordWrong)
	}

	accessToken, refreshToken, err := common.GenToken(id, name)
	if err != nil {
		return
	}

	user = &common.User{
		UserID:       id,
		UserName:     name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return user, err
}

const secret = "LICSLAN"

// encryptPassword 对密码进行加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

func SignUp(p *common.SignUpData) (error error) {
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

	sqlStr := `SELECT * FROM t_user WHERE username = ?`

	var us []common.User
	rows, err := Db.Dbs.Query(sqlStr, username)
	if err != nil {
		return err
	}
	for rows.Next() {
		var u common.User
		// sqlx 提供了便捷方法可以将查询结果直接扫描到结构体
		err = rows.Scan(&u)
		if err != nil {
			return err
		}
		us = append(us, u)
	}

	if cap(us) > 0 {
		return errors.New(common.ErrorUserExit)
	}
	return
}

// SaveUser 注册业务-向数据库中插入一条新的用户
func SaveUser(user common.User) (error error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	// 执行SQL语句入库
	sqlStr := `INSERT INTO t_user(user_id,username,password,email,gender) VALUES(?,?,?,?,?)`
	_, err := Db.Dbs.Exec(sqlStr, user.UserID, user.UserName, user.Password, user.Email, user.Gender)
	return err
}
