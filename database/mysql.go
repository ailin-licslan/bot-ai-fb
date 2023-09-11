package database

import (
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

func InitDB() (err error) {

	//Dbs, err := sql.Open("mysql", str) 注意这里不要冒号赋值!!!
	Dbs, err = sql.Open("mysql", "root:123456@tcp(192.168.0.161:3306)/licslan")
	if err != nil {
		return err
	}
	//defer Dbs.Close()
	//max connect num
	Dbs.SetMaxOpenConns(30)
	//max idle num
	Dbs.SetMaxIdleConns(2)
	errTest := Dbs.Ping()
	if errTest != nil {
		fmt.Print(errTest)
	}
	fmt.Println("database connected success!")
	return nil
}
