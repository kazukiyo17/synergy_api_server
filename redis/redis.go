package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"log"
	"time"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data string, expire int) {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, data)
	if err != nil {
		log.Printf("redis set error: %v", err)
		return
	}
	// expire 以天为单位
	_, err = conn.Do("EXPIRE", key, expire*24*3600)
	if err != nil {
		log.Printf("redis set expire error: %v", err)
		return
	}
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) string {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Printf("redis get error: %v", err)
		return ""
	}
	return string(reply)
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}
