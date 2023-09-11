package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var Con redis.Conn

func InitRedis() {

	Con, _ = redis.Dial("tcp", "192.168.0.161:6380")
	//defer Con.Close()
	//PASSWORD IF YOU HAVE ONE
	//dial.Do("Auth","123456")
	fmt.Print(Con)

}
