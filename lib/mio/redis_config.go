package mio

import (
	"game-im/config"
)

var (
	DazaRedisSlave                   = NewRedisClient(config.REDIS_DAZA_SLAVE)
	SessionRedis                     = NewRedisClient(config.REDIS_SESSION_MASTER)
	ProfileaRedis      *RedisCluster = NewRedisCluster(config.REDIS_PROFILEA_SLAVE_TPL, config.REDIS_PROFILEA_SLAVE_BASE, config.REDIS_PROFILEA_SLAVE_STEP, nil)
	GameInfoRedisSlave *RedisClient  = NewRedisClient(config.REDIS_GAMEINFO_SLAVE)
)
