package global

import (
	"github.com/go-redis/redis"
	"time"
)
var(RedisDB *redis.Client)
const Duration = time.Minute * 30
func SetupRedisDb() error{
	RedisDB = redis.NewClient(&redis.Options{
		Addr: ":6379",
		Password: "",
		DB: 0,
	})
	_,err := RedisDB.Ping().Result()
	if err != nil{
		return err
	}
	return nil
}