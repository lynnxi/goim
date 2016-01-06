package mio

import (
	"fmt"
	"strconv"
	"strings"
)

type RedisCluster struct {
	redisClients map[string]*RedisClient
	hostTpl      string
	step         int
	base         int
	shardFunc    func(key string, hostTpl string, base int, step int) (host string)
}

func NewRedisCluster(hostTpl string, base int, step int, shardFunc func(key string, hostTpl string, base int, step int) (host string)) (cluster *RedisCluster) {
	if shardFunc == nil {
		shardFunc = ExplodeModShardFuc
	}
	cluster = &RedisCluster{
		redisClients: make(map[string]*RedisClient),
		hostTpl:      hostTpl,
		shardFunc:    shardFunc,
		base:         base,
		step:         step,
	}

	return
}

func (cluster *RedisCluster) getRedis(key string) (client *RedisClient) {
	host := cluster.shardFunc(key, cluster.hostTpl, cluster.base, cluster.step)
	client, ok := cluster.redisClients[host]
	if !ok {
		client = NewRedisClient(host)
		cluster.redisClients[host] = client
	}

	return
}

func (cluster *RedisCluster) Get(key string) (value interface{}, err error) {
	value, err = cluster.getRedis(key).Get(key)
	return
}

func (cluster *RedisCluster) Hset(key string, field string, value string) (err error) {
	_, err = cluster.getRedis(key).Hset(key, field, value)
	return
}

func (cluster *RedisCluster) Zrange(key string, start int, end int, withscores bool) (value interface{}, err error) {
	value, err = cluster.getRedis(key).Zrange(key, start, end, withscores)

	return
}

func ExplodeModShardFuc(key string, hostTpl string, base int, step int) (host string) {
	parts := strings.Split(key, ":")
	id, _ := strconv.Atoi(parts[1])
	n := id % base
	var c int
	if step == 0 {
		c = 0
	} else {
		c = id / step
	}
	host = fmt.Sprintf(hostTpl, c, n, n)
	return
}
