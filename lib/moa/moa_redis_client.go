package moa

import (
	//"encoding/json"
	"errors"
	"game-im/lib/hash"
	"game-im/lib/redigo/redis"
	//"game-im/lib/stdlog"
	"game-im/lib/vmihailenco/msgpack"
	"math/rand"
	"os"
	"time"
)

type MoaRedisClient struct {
	hosts       []string
	hash        int
	hashLocator *hash.KetamaHashLocator
	serviceUri  string
	clientCache map[string]*redis.Pool
	time        int64
}

func NewMoaRedisClient(hosts []string, serviceUri string, hash int) (client *MoaRedisClient) {
	client = &MoaRedisClient{hosts: hosts,
		serviceUri:  serviceUri,
		clientCache: make(map[string]*redis.Pool),
		hash:        hash,
		time:        time.Now().Unix()}
	outLogMoa, _ := os.OpenFile("/home/deploy/log/server_moa.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	moaLogger.SetOutput(outLogMoa)

	return
}

func (client *MoaRedisClient) IsExpired() bool {
	var t int64 = 5
	now := time.Now().Unix()
	return now-t > client.time
}

func (client *MoaRedisClient) Invoke(method string, args []interface{}, hashKey string) (result interface{}, err error) {
	// 构建将要发送到moa的指令
	cmd := map[string]interface{}{
		"id":     GenerateUid(),
		"action": client.serviceUri,
		"params": map[string]interface{}{
			"m":    method,
			"args": args,
		},
	}

	_cmdjson, _ := msgpack.Marshal(cmd)
	cmdjson := string(_cmdjson)
	conn := client.getRedis(hashKey).Get()
	defer conn.Close()
	respjson, err := conn.Do("GET", cmdjson)

	if err != nil {
		return
	}
	//	moaLogger.Println("moa resp : " + string(respjson.([]byte)))

	var rtn map[string]interface{}
	if err = msgpack.Unmarshal(respjson.([]byte), &rtn); err != nil {

	}
	if int(rtn["ec"].(int64)) != 0 {
		err = errors.New("moa get error : " + string(respjson.([]byte)))
	}
	result = rtn["result"]

	return
}

func (client *MoaRedisClient) getRedis(hashKey string) (r *redis.Pool) {
	var x int
	var host string
	var n = len(client.hosts)

	if n == 0 {

	}

	if n > 1 {
		if client.hash == KETAMA {
			if client.hashLocator == nil {
				client.hashLocator = hash.NewKetamaHashLocator(client.hosts)
			}
			host = client.hashLocator.GetNodeByKey([]byte(hashKey))
		} else {
			x = rand.Intn(n - 1)
			host = client.hosts[x]
		}
	} else {
		host = client.hosts[0]
	}

	//moaLogger.Println("use moa host : " + host + " hash key : " + hashKey)
	r, ok := client.clientCache[host]
	if !ok {
		r = &redis.Pool{
			MaxIdle:     2000,
			IdleTimeout: 60 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host)
				if err != nil {
					return nil, err
				}

				return c, err
			},
		}
		client.clientCache[host] = r
	}

	return
}
