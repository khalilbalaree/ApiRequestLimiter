package conf

type Redis struct {
	Host string `json:"host"`
}

type Limiter struct {
	MaxPermits string `json:"maxPermits"`
	Rate       string `json:"rate"`
}

var (
	redis   Redis
	limiter Limiter
)

func init() {
	//TODO: load file
}

func GetRedis() Redis {
	return redis
}

func GetLimiter() Limiter {
	return limiter
}
