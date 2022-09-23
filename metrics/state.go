package metrics

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Metrics struct {
	name              string
	redisCli          *redis.Client
	defaultExpiration time.Duration
}

//127.0.0.1:6379> HINCRBY topic succ 1
