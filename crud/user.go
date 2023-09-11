package crud

import (
	"bot-ai-fb/model"
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

func InsertUser(db *sql.DB, user model.User) {

	db.Begin()
	add, err := db.Prepare("insert into user (name,age) values (?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = add.Exec(user.Name, user.Age)
	if err != nil {
		log.Fatal(err)
	}
	add.Close()
	defer db.Close()

}

func DeleteUser(db *sql.DB, id int) {

	db.Begin()
	del, err := db.Prepare("delete from user where id =?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = del.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	del.Close()
	defer db.Close()

}

func RedisAdd(con redis.Conn, name string) {

	//set
	con.Do("Hset", "u", "name", name)

	//get
	res, err := redis.String(con.Do("Hget", "u", "name"))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
