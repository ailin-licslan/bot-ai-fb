package redis

import (
	settings "bot-ai-fb/setting"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var Con redis.Conn

func InitRedis(cfg *settings.RedisConfig) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	//"192.168.0.161:6380"
	Con, _ = redis.Dial("tcp", addr)
	//defer Con.Close()
	//PASSWORD IF YOU HAVE ONE
	//dial.Do("Auth","123456")
	fmt.Print(Con)

}

func Close() {
	_ = Con.Close()
}
