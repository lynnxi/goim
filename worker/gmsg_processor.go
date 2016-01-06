package main

import (
	"encoding/json"
	"game-im/config"
	. "game-im/lib/mio"
	//"game-im/lib/stdlog"
	"game-im/lib/mcrypt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

var (
	redis6080    = NewRedisClient(config.REDIS_DAZA_SLAVE)
	gwQueueRedis = map[string]*RedisClient{}
	wait         sync.WaitGroup
	Moa          *MoaClient = NewMoaClient()
)

func PopGmsg() {
	defer func() {
		errorLogger.Println("pop msg shut down...")
	}()

	for _, host := range config.ImqRedisConfig {
		wait.Add(1)
		go popGmsg(host)
	}

	wait.Wait()
}

func popGmsg(host string) {
	defer func() {
		if err := recover(); err != nil { //子协程崩溃貌似catch不住
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
		}
		wait.Done()
	}()
	var msg *Msg
	var msgs []*Msg
	var reqPack *Gpack
	var t, text, fr, to, appid, gw_host, key string
	var imj, _server, union map[string]interface{}
	var v interface{}
	var err error
	var _qredis *RedisClient
	var ok bool
	var n = 0
	var now = time.Now().Unix()

	redis := NewRedisClient(host)
	for {

		n++
		_now := time.Now().Unix()
		if now != _now {
			infoLogger.Println("deal with " + config.KEY_IMQ_V2 + " " + strconv.Itoa(n) + " per second")
			n = 0
			now = _now
		}

		v, err = redis.Blpop(config.KEY_IMQ_V2, 0)
		if err != nil {
			errorLogger.Println("BLPOP ERR : ", err)
			continue
		}

		imj = nil
		if err = json.Unmarshal(v.([]byte), &imj); err != nil {
			errorLogger.Println("RedisPool: json decoding failed: %s", err)
			errorLogger.Println(string(v.([]byte)))
			continue
		}

		if source, ok := imj["source"]; ok && source == "game" { //游戏发到群里的消息 不再发回去
			continue
		}

		//查询绑定关系

		key = "groupid:" + imj["to"].(string) + ":union"
		//key = "groupid:33567928:union"

		union, err = redis6080.HgetAll(key)

		if union == nil {
			continue
		}

		if err != nil {
			errorLogger.Println("group binding get error ", err)
			continue
		}

		// group, err = Moa.GetGroupProfile(imj["to"].(string), imj["fr"].(string))

		// if err != nil {
		// 	errorLogger.Println("group binding get error ", err)
		// 	continue
		// }
		// if group["game_union"] == nil {
		// 	continue
		// }
		// union = group["game_union"].(map[interface{}]interface{})

		//拼装数据
		text = imj["text"].(string)
		fr = imj["fr"].(string)
		to = union["unionid"].(string)
		appid = union["appid"].(string)
		switch imj["type"].(type) {
		case float64:
			t = strconv.Itoa(int(imj["type"].(float64)))
			break
		case string:
			t = imj["type"].(string)
			j, _ := json.Marshal(imj)
			errorLogger.Println("imj get string type ", string(j))
			break
		default:
			t = "1"
			j, _ := json.Marshal(imj)
			errorLogger.Println("imj get unknown type ", string(j))
			break
		}

		//处理消息
		switch t {
		case "1":
			break
		case "2":
			text = "给您发来一张图片，请到陌陌查看"
			break
		case "3":
			text = "给您发来一段语音，请到陌陌查看"
			break
		case "4":
			text = "给您发来ta的位置，请到陌陌查看"
			break
		case "5":
			text = "给您发来一个陌陌表情"
			break
		case "7":
			text = "给您发来一段视频，请到陌陌查看"
			break
		default:
			text = "当前版本不支持这个信息，请到陌陌查看"
			break
		}

		if len(text) > 130 {
			text = "我讲了一大段故事，你可以在陌陌慢慢看"
		}

		//取名字
		name, err := ProfileaRedis.Get("user:" + fr + ":name")
		if err != nil {
			errorLogger.Println("GET USER NAME ERR : ", err)
			name = "momo"
		}

		if appid == "ex_mmzb_wFaBpbG" || appid == "ex_mmzbtest_4QGQJKy" {
			fr, _ = mcrypt.EncryptV1([]byte(fr), config.SALT)
		} else {
			fr, _ = mcrypt.EncryptV2([]byte(fr), []byte(strconv.Itoa(config.SALT)+appid))
		}

		msg = NewMsg(fr, to, text, t, strconv.Itoa(int(time.Now().Unix())), 1, string(name.([]byte)))
		msgs = append(msgs, msg)
		reqPack = NewReqGpack("ugmsg", map[string]interface{}{"count": 1, "msgs": msgs})
		reqPack.SetFlag(appid) //标记数据包sessionid

		//push到gateway
		if true {
			//查询appid连在哪个机器
			v, err := SessionRedis.Hget("session_map", appid)
			if err != nil {
				errorLogger.Println("app session map not found ", appid)
				continue
			}
			if v == nil {
				errorLogger.Println("app session map not found ", appid)
				continue
			}
			if err = json.Unmarshal(v.([]byte), &_server); err != nil {
				errorLogger.Println("RedisPool: json decoding failed: %s", err)
				errorLogger.Println(string(v.([]byte)))
				continue
			}

			// ips := make([]string, 0, len(list))
			// for k := range list {
			// 	ips = append(ips, k)
			// }

			// //按发送目标hash
			// i := int(crc32.ChecksumIEEE([]byte(to)))
			// i = i % len(ips)
			// ip := ips[i]
			gw_host = _server["local_host"].(string) + ":" + _server["local_port"].(string)
			_qredis, ok = gwQueueRedis[gw_host]
			if !ok {
				_qredis = NewRedisClient(gw_host)
				gwQueueRedis[gw_host] = _qredis
			}
			key = config.KEY_GW_IN
			_qredis.Rpush(key, string(reqPack.ToBytes()))
			msgs = nil //清空数组
		}

		defer func() {
			if err := recover(); err != nil { //子协程崩溃貌似catch不住
				panicLogger.Println("panic : ", err)
				panicLogger.Println("stack : ", string(debug.Stack()))
			}
		}()
	}

}
