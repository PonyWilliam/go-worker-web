package cache

import (
	"encoding/json"
	"github.com/PonyWilliam/go-WorkWeb/global"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro/v2/util/log"
)
func GetGlobalCache(key string)(rsp interface{},err error){
	val,err := global.RedisDB.Get(key).Result()
	if err == redis.Nil || err != nil{
		return nil,err
	}
	return val,nil
}
func SetGlobalCache(key string,res interface{}) error {
	content,err := json.Marshal(res)
	if err != nil{
		log.Fatal(err)
	}
	err = global.RedisDB.Set(key,content,global.Duration).Err()
	if err != nil{
		log.Fatal(err)
	}
	return nil
}
func DelCache(key string){
	global.RedisDB.Del(key)
}