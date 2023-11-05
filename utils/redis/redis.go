package redis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	//"github.com/go-redis/redis/v8"
	"github.com/kazukiyo17/fake_buddha_server/common/conf"
)

// Exists 判断key是否存在
func Exists(key string) bool {
	conn := conf.C.RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := conf.C.RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Set(key string, data interface{}, time int) error {
	conn := conf.C.RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}
