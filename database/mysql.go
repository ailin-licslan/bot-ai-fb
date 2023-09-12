package database

import (
	settings "bot-ai-fb/setting"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() {

	dsn := "root:123456@tcp(192.168.0.161:3306)/licslan"

	connect, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Print(err)
	}

	errTest := connect.Ping()
	if errTest != nil {
		fmt.Print(errTest)
	}

	fmt.Println("database connected success!")

}

var Dbs *sql.DB

func InitDB(cfg *settings.MySQLConfig) (err error) {

	//Dbs, err := sql.Open("mysql", str) 注意这里不要冒号赋值!!!

	//"root:123456@tcp(192.168.0.161:3306)/licslan"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB)

	Dbs, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	//defer Dbs.Close()
	//max connect num
	Dbs.SetMaxOpenConns(cfg.MaxOpenConns)
	//max idle num
	Dbs.SetMaxIdleConns(cfg.MaxIdleConns)
	errTest := Dbs.Ping()
	if errTest != nil {
		fmt.Print(errTest)
	}
	fmt.Println("database connected success!")
	return nil
}

func Close() {
	_ = Dbs.Close()
}
