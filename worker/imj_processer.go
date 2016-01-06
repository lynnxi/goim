package main

// import (
// 	"encoding/json"
// 	"game-im/config"
// 	"net/http"
// 	//"game-im/lib/stdlog"
// 	"hash/crc32"
// 	"io/ioutil"
// 	"runtime/debug"
// 	"strconv"
// 	"time"
// )

// var (
// 	dazaRedis  = NewRedisClient(config.REDIS_DAZA_SLAVE)
// 	queueRedis = map[string]*RedisClient{}
// 	sessionMap map[string]map[string]interface{} //这两个变量线程不安全 需要加锁
// 	serverList = [...]string{"10.83.60.90", "10.83.60.92"}
// )

// func GetSessionMap() {

// 	var data map[string]interface{}
// 	var _sessionMap = make(map[string]map[string]interface{})

// 	for ip, _ := range queueRedis { //删除不存在得机器redis链接
// 		in := false
// 		for _, v := range serverList {
// 			if v == ip {
// 				in = true
// 			}
// 		}
// 		if !in {
// 			delete(queueRedis, ip)
// 		}
// 	}

// 	for _, ip := range serverList {
// 		_, ok := queueRedis[ip]
// 		if !ok {
// 			queueRedis[ip] = NewRedisClient(ip + ":1603")
// 		}
// 		res, err := http.Get("http://" + ip + ":6060/info")
// 		if err != nil {
// 			errorLogger.Println("session get error, ", err, ip)
// 			continue
// 		}
// 		result, err := ioutil.ReadAll(res.Body)
// 		res.Body.Close()
// 		if err != nil {
// 			errorLogger.Println("session get error, ", err, ip)
// 			continue
// 		}
// 		json.Unmarshal(result, &data)
// 		for appid, time := range data["session"].(map[string]interface{}) {
// 			_, ok := _sessionMap[appid]
// 			if !ok {
// 				_sessionMap[appid] = make(map[string]interface{})
// 			}
// 			_sessionMap[appid][ip] = time
// 		}
// 	}
// 	sessionMap = _sessionMap
// }

// func PopImj() {
// 	defer func() {
// 		errorLogger.Println("pop msg shut down...")
// 	}()

// 	for _, host := range config.ImqRedisConfig {
// 		go popImj(host)
// 	}

// 	for {
// 		time.Sleep(1 * time.Second)
// 	}

// }

// func popImj(host string) {
// 	defer func() {
// 		if err := recover(); err != nil { //子协程崩溃貌似catch不住
// 			panicLogger.Println("panic : ", err)
// 			panicLogger.Println("stack : ", string(debug.Stack()))
// 		}
// 	}()
// 	var t string
// 	var n int
// 	var v interface{}
// 	var err error
// 	var imj map[string]interface{}
// 	var bindingInfo map[string]interface{}
// 	var _qredis *RedisClient
// 	redis := NewRedisClient(host)
// 	for {
// 		v, err = redis.Blpop(config.KEY_IMQ_OUT, 0)
// 		if err != nil {
// 			errorLogger.Println("BLPOP ERR : ", err)
// 			n++
// 			continue
// 		}
// 		if v == nil {
// 			n++
// 		} else {
// 			n = 0

// 			imj = nil
// 			if err = json.Unmarshal(v.([]byte), &imj); err != nil {
// 				errorLogger.Println("RedisPool: json decoding failed: %s", err)
// 				errorLogger.Println(string(v.([]byte)))
// 				continue
// 			}

// 			if source, ok := imj["source"]; ok && source == "game" {
// 				//errorLogger.Println("game msg : ", string(v.([]byte)))
// 				continue
// 			}

// 			// if imj["fr"] == "5125453" {
// 			// 	j, _ := json.Marshal(imj)
// 			// 	errorLogger.Println("aguai string : ", string(v.([]byte)))
// 			// 	errorLogger.Println("a guai msg : ", string(j))
// 			// }

// 			//查绑定关系

// 			// if imj["to"].(string) == "18484752" {
// 			// 	infoLogger.Println(imj["text"])
// 			// }
// 			key := "group:" + imj["to"].(string) + ":binding"
// 			v, err = dazaRedis.Get(key)

// 			if v == nil {
// 				//stdlog.Println("group binding get nil ")
// 				continue
// 			}

// 			if err != nil {
// 				errorLogger.Println("group binding get error ", err)
// 				continue
// 			}
// 			if err = json.Unmarshal(v.([]byte), &bindingInfo); err != nil {
// 				errorLogger.Println("RedisPool: json decoding failed: %s", err)
// 				errorLogger.Println(string(v.([]byte)))
// 				continue
// 			}

// 			//查session
// 			// key = "session_map"
// 			// v, err = dazaRedis.Hget(key, bindingInfo["appid"].(string)) //这里可以优化 例如zookeeper
// 			// if v == nil {
// 			// 	noticeLogger.Println("session not exist : ", bindingInfo["appid"])
// 			// 	continue
// 			// }

// 			// if err != nil {
// 			// 	errorLogger.Println("session get error ", err)
// 			// 	continue
// 			// }
// 			list, ok := sessionMap[bindingInfo["appid"].(string)]
// 			if !ok {
// 				//noticeLogger.Println("session not exist : ", bindingInfo["appid"])
// 				continue
// 			}

// 			ips := make([]string, 0, len(list))
// 			for k := range list {
// 				ips = append(ips, k)
// 			}
// 			i := int(crc32.ChecksumIEEE([]byte(imj["to"].(string))))
// 			i = i % len(ips)
// 			ip := ips[i]
// 			_qredis = queueRedis[ip]
// 			//写队列
// 			key = "appid:" + bindingInfo["appid"].(string) + ":gpack_queue"

// 			switch imj["type"].(type) {
// 			case float64:
// 				t = strconv.Itoa(int(imj["type"].(float64)))
// 				break
// 			case string:
// 				t = imj["type"].(string)
// 				j, _ := json.Marshal(imj)
// 				errorLogger.Println("imj get string type ", string(j))
// 				break
// 			default:
// 				t = "1"
// 				j, _ := json.Marshal(imj)
// 				errorLogger.Println("imj get unknown type ", string(j))
// 				break

// 			}

// 			gpackjson, _ := json.Marshal(map[string]interface{}{
// 				"cmd":   "msg",
// 				"appid": bindingInfo["appid"],
// 				"body": map[string]interface{}{
// 					"from": imj["fr"],
// 					"to":   bindingInfo["allyid"].(string),
// 					"text": imj["text"],
// 					"type": t}})

// 			if true {
// 				_qredis.Rpush(key, string(gpackjson))

// 			}
// 			defer func() {
// 				if err := recover(); err != nil { //子协程崩溃貌似catch不住
// 					panicLogger.Println("panic : ", err)
// 					panicLogger.Println("stack : ", string(debug.Stack()))
// 				}
// 			}()
// 		}
// 	}
// 	if n >= 5 {
// 		noticeLogger.Println("blpop get nil 5th %s", n, " ", host)
// 		time.Sleep(time.Second)
// 		n = 0
// 	}

// }
