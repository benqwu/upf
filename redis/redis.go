package redis

import (
	"N4test/conf"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var redisConn *redis.Pool

// GetRedisConn 获取 Redis 客户端连接
func GetRedisConn() *redis.Pool {
	return redisConn
}

// Setup 创建 Redis 连接
func Setup() error {
	redisConn = &redis.Pool{
		MaxIdle:     conf.RedisConf.MaxIdle,
		MaxActive:   conf.RedisConf.MaxActive,
		IdleTimeout: conf.RedisConf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisConf.Host+":"+strconv.Itoa(conf.RedisConf.Port))
			if err != nil {
				return nil, err
			}
			// 验证密码
			if conf.RedisConf.Password != "" {
				if _, err := c.Do("AUTH", conf.RedisConf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 选择数据库
			if _, err := c.Do("SELECT", conf.RedisConf.DBNum); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// Set 方法
func Set(key string, data string, seconds int) error {
	conn := GetRedisConn().Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return err
	}
	return nil
}

// Exists 方法
func Exists(key string) bool {
	conn := GetRedisConn().Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// Get 方法
func Get(key string) (string, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	if err == redis.ErrNil {
		return "", nil
	}
	return reply, nil
}

// Del 方法
func Del(key string) (bool, error) {
	conn := GetRedisConn().Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// DelLike 方法
func DelLike(key string) error {
	conn := GetRedisConn().Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err := Del(key)
		if err != nil {
			return err
		}
	}
	return nil
}
