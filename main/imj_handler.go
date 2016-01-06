package main

import (
	"encoding/json"
	"game-im/config"
	//"game-im/lib/stdlog"
	"time"
)

func (server *Server) PopImj() {
	defer func() {
		errorLogger.Println("pop msg shut down...")
	}()

	queueRedis := map[string]*RedisClient{}

	for _, host := range config.ImqRedisConfig {
		queueRedis[host] = NewRedisClient(host)
	}
	dazaRedis := NewRedisClient(config.REDIS_DAZA_SLAVE)

	var n int
	var v interface{}
	var err error
	for {

		for _, redis := range queueRedis {
			v, err = redis.Lpop(config.KEY_IMQ)
			if err != nil {

			}
			if v == nil {
				n++
			} else {
				n = 0
				var imj *Imj = new(Imj)
				if err = json.Unmarshal(v.([]byte), imj); err != nil {
					errorLogger.Println("RedisPool: json decoding failed: %s", err)
					errorLogger.Println(string(v.([]byte)))
					continue
				}

				key := "group:" + imj.To + ":binding"
				v, err = dazaRedis.Get(key)

				if v == nil {
					//stdlog.Println("group binding get nil ")
					continue
				}

				if err != nil {
					errorLogger.Println("group binding get error ", err)
					continue
				}
				var data map[string]interface{}
				if err = json.Unmarshal(v.([]byte), &data); err != nil {
					errorLogger.Println("RedisPool: json decoding failed: %s", err)
					errorLogger.Println(string(v.([]byte)))
					continue
				}

				s, exist := server.SessionMap[data["appid"].(string)]
				if exist {
					//满时 关闭这个链接
					if len(s.ReqPackBufferChannel) >= config.IMJ_CHANNEL_BUFF_SIZE {
						noticeLogger.Println("chan buff full")

					} else {
						str, err := json.Marshal(map[string]string{"from": imj.Fr,
							"to":   data["allyid"].(string),
							"body": imj.Text,
							"type": imj.Type})
						if err != nil {

						}

						gpack := NewGpack(data["appid"].(string), "msg_send", string(str))

						s.ReqPackBufferChannel <- gpack
					}
				} else {
					noticeLogger.Println("session not exist : %s", data["appid"])
				}
			}
		}

		if n >= 200 {
			noticeLogger.Println("lpop get nil 20th %s", n)
			time.Sleep(time.Second)
			n = 0
		}

	}
}
