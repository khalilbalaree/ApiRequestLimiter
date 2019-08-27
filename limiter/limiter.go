package limiter

import (
	"errors"
	// "git.code.oa.com/cloud_industry/boss/job/conf"
	"sync"
	"time"

	"github.com/ApiRequestLimiter/data"
	"github.com/gomodule/redigo/redis"
)

var limiterAgent *LimiterAgent
var once sync.Once

type LimiterValue struct {
	MaxPermits string
	Rate       string
}

type LimiterAgent struct {
	spec   LimiterValue
	pool   *redis.Pool
	script *redis.Script
}

func GLimiterAgent() *LimiterAgent {
	once.Do(func() {
		limiterAgent = NewLimiterAgent()
	})
	return limiterAgent
}

func NewLimiterAgent() *LimiterAgent {
	if data.RedisPool() == nil {
		panic(errors.New("pool error"))
		return nil
	}

	limiterAgent := &LimiterAgent{
		spec: LimiterValue{
			//配置文件中读取
			MaxPermits: conf.GetLimiter().MaxPermits,
			Rate:       conf.GetLimiter().Rate,
		},
		pool:   data.RedisPool(),
		script: redis.NewScript(1, luaText),
	}
	return limiterAgent
}

//入口
//先获得锁，再处理
func (l *LimiterAgent) HandleRequest(user string, numRequest int64) (bool, error) {
	lockKey := getLimiterLockKey(user)

	for {
		//获得锁
		lock, err := l.limiterGetLock(lockKey)
		if err != nil {
			return false, err
		}
		if lock {
			break
		}
	}

	result, err := l.DoLimit(user, numRequest, time.Now().UnixNano())
	if err != nil {
		return false, err
	}

	//解锁
	if err := l.limiterUnLock(lockKey); err != nil {
		return false, err
	}

	return result, nil
}

func (l *LimiterAgent) DoLimit(user string, numRequest int64, currNanoSec int64) (bool, error) {
	conn := l.pool.Get()

	defer conn.Close()

	key := getLimiterKey(user)
	value, err := l.script.Do(conn, key, currNanoSec, numRequest, l.spec.MaxPermits, l.spec.Rate)
	if err != nil {
		return false, err
	}
	if value != nil {
		return true, nil
	}
	return false, nil

}

//lock
func (l *LimiterAgent) limiterGetLock(lockKey string) (bool, error) {
	conn := l.pool.Get()
	_, err := redis.String(conn.Do("SET", lockKey, 1, "EX", 1, "NX"))
	if err == redis.ErrNil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

//unlock
func (l *LimiterAgent) limiterUnLock(lockKey string) error {
	conn := l.pool.Get()
	_, err := conn.Do("DEL", lockKey)
	if err != nil {
		return err
	}
	return nil
}
