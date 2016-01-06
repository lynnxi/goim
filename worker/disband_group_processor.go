package main

import (
	"encoding/json"
	"game-im/config"
	. "game-im/lib/mio"
	//"game-im/lib/stdlog"
	"game-im/lib/mcrypt"
	"runtime/debug"
	"strconv"
)

func DisbandGroup() {
	defer func() {
		errorLogger.Println("disband group shut down...")
	}()

	for _, host := range config.ImqRedisConfig {
		wait.Add(1)
		go disbandGroup(host)
	}

	wait.Wait()
}

func disbandGroup(host string) {
	defer func() {
		if err := recover(); err != nil { //子协程崩溃貌似catch不住
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
		}
		wait.Done()
	}()

	var reqPack *Gpack
	var appid, gw_host string
	var unionid, operator string
	var imj, _server map[string]interface{}
	var v interface{}
	var err error
	var _qredis *RedisClient
	var ok bool
	redis := NewRedisClient(host)
	for {
		v, err = redis.Brpop("disband_game_group_notify", 0)
		if err != nil {
			errorLogger.Println("BRPOP ERR : ", err)
			continue
		}

		imj = nil
		if err = json.Unmarshal(v.([]byte), &imj); err != nil {
			errorLogger.Println("RedisPool: json decoding failed: %s", err)
			errorLogger.Println(string(v.([]byte)))
			continue
		}

		//拼装数据
		appid = imj["appid"].(string)
		unionid = imj["unionid"].(string)
		operator = imj["operator"].(string)

		//查询appid连在哪个机器
		v, err := SessionRedis.Hget("session_map", appid)
		if err != nil {
			errorLogger.Println("app session map not found ", appid)
			continue
		}
		if v == nil {
			continue
		}
		if err = json.Unmarshal(v.([]byte), &_server); err != nil {
			errorLogger.Println("RedisPool: json decoding failed: %s", err)
			errorLogger.Println(string(v.([]byte)))
			continue
		}

		gw_host = _server["local_host"].(string) + ":" + _server["local_port"].(string)
		_qredis, ok = gwQueueRedis[gw_host]
		if !ok {
			_qredis = NewRedisClient(gw_host)
			gwQueueRedis[gw_host] = _qredis
		}

		if appid == "ex_mmzb_wFaBpbG" || appid == "ex_mmzbtest_4QGQJKy" {
			operator, _ = mcrypt.EncryptV1([]byte(operator), config.SALT)
		} else {
			operator, _ = mcrypt.EncryptV2([]byte(operator), []byte(strconv.Itoa(config.SALT)+appid))
		}

		reqPack = NewReqGpack("disbandGroup", map[string]interface{}{"appid": appid, "unionid": unionid, "operator": operator})
		reqPack.SetFlag(appid) //标记数据包sessionid

		//push到gateway
		key := config.KEY_GW_IN
		_qredis.Rpush(key, string(reqPack.ToBytes()))

		defer func() {
			if err := recover(); err != nil { //子协程崩溃貌似catch不住
				panicLogger.Println("panic : ", err)
				panicLogger.Println("stack : ", string(debug.Stack()))
			}
		}()
	}
}
