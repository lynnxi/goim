package main

import (
	"encoding/json"
	"errors"
	"game-im/config"
	. "game-im/lib/mio"
	//"io"
	"runtime/debug"
	//"strconv"
)

type Logic struct {
	local_host              string
	local_port              string
	logicQueueRedis         *RedisClient
	InputPackBufferChannel  chan *Gpack
	OutputPackBufferChannel chan *Gpack
}

func newLogic(local_host string, local_port string) (logic *Logic) {
	logic = &Logic{
		local_host:              local_host,
		local_port:              local_port,
		logicQueueRedis:         NewRedisClient(local_host + ":" + local_port),
		InputPackBufferChannel:  make(chan *Gpack, config.CHANNEL_BUFF_SIZE),
		OutputPackBufferChannel: make(chan *Gpack, config.CHANNEL_BUFF_SIZE),
	}
	return logic
}

func (logic *Logic) Start() {
	go logic.PopGpack()
	go logic.PushGpack()
	logic.HandleGpack()

}

func (logic *Logic) PopGpack() {
	var gpack *Gpack
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
		}
		noticeLogger.Println("PopGpack close...")

	}()
	for {
		v, err := logic.logicQueueRedis.Blpop(config.KEY_LOGIC_IN, 0)
		if err != nil {
			errorLogger.Println("BLPOP ERR : ", err)
			continue
		}
		gpack = nil
		if err = json.Unmarshal(v.([]byte), &gpack); err != nil {
			errorLogger.Println("RedisPool: json decoding failed: %s", err)
			errorLogger.Println(string(v.([]byte)))
			continue
		}
		logic.InputPackBufferChannel <- gpack
	}
}

func (logic *Logic) PushGpack() {
	var gpqak *Gpack
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
		}
		noticeLogger.Println("PushGpack close...")

	}()
	for {
		gpqak = <-logic.OutputPackBufferChannel
		if gpqak == nil {
			return
		}

		host, err := logic.getGwServer(gpqak.GetFlag())
		if err != nil {
			errorLogger.Println("get gw error ", err.Error())
			continue
		}
		redis, ok := gwQueueRedis[host]
		if !ok {
			redis = NewRedisClient(host)
			gwQueueRedis[host] = redis
		}
		redis.Rpush(config.KEY_GW_IN, string(gpqak.ToBytes()))
	}
}

func (logic *Logic) getGwServer(sessionid string) (host string, err error) {
	//查询appid连在哪个机器
	v, err := SessionRedis.Hget("session_map", sessionid)
	if err != nil {
		errorLogger.Println("app session map not found ", sessionid)
		return
	}
	if v == nil {
		errorLogger.Println("app session map not found ", sessionid)
		err = errors.New("app session map not found " + sessionid)
		return
	}
	var _server map[string]interface{}
	if err = json.Unmarshal(v.([]byte), &_server); err != nil {
		errorLogger.Println("RedisPool: json decoding failed: %s", err)
		errorLogger.Println(string(v.([]byte)))
		return
	}
	host = _server["local_host"].(string) + ":" + _server["local_port"].(string)
	// ips := make([]string, 0, len(list))
	// for k := range list {
	// 	ips = append(ips, k)
	// }

	// //按发送目标hash
	// i := int(crc32.ChecksumIEEE([]byte(to)))
	// i = i % len(ips)
	// ip := ips[i]

	//host = config.REDIS_GW_1
	return
}

func (logic *Logic) HandleGpack() {
	var req, ret *Gpack

	for {
		req = <-logic.InputPackBufferChannel
		if req == nil { //channel 被关闭
			return
		}
		ret = logic.exec(req)
		logic.OutputPackBufferChannel <- ret
	}
}

func (logic *Logic) exec(req *Gpack) (res *Gpack) {
	var err error
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
			res = GetErrGpack("error", req.Sid)
			noticeLogger.Println("exec err...")
			return
		}
	}()
	f, ok := Action2Handler[req.Cmd]
	if ok {
		res, err = f(req)
		if err != nil {
			res = GetErrGpack(err.Error(), req.Sid)
		}

	} else {
		res = GetErrGpack(err.Error(), req.Sid)
	}
	res.SetFlag(req.GetFlag())

	return
}
