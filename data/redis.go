package data

import (
	"github.com/ApiRequestLimiter/conf"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

var (
	redisPool *redis.Pool
)

func init() {
	var err error
	redisPool, err = newLocalRedis()
	if err != nil {
		logrus.Fatal("redis init fail, ", err)
		return
	}
}

//单点redis
func newLocalRedis() (*redis.Pool, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.GetRedis().Host)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
	return pool, nil
}

func RedisPool() *redis.Pool {
	return redisPool
}
