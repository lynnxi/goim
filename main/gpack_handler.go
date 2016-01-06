package main

import (
	//"game-im/config"
	"encoding/json"
	"errors"
	//"game-im/lib/stdlog"
	//"time"
)

var Action2Handler = map[string]func(req *Gpack, server *Server) (res *Gpack, err error){
	"auth": GpackAuthHandler,
	"info": GpackInfoHandler,
}

func GpackInfoHandler(req *Gpack, server *Server) (res *Gpack, err error) {
	infoMap := make(map[string]map[string]int)
	for appid, session := range server.SessionMap {
		infoMap[appid] = map[string]int{
			"read":          session.conn.stats["read_qpack"],
			"write":         session.conn.stats["write_qpack"],
			"req_buf_chan":  len(session.ReqPackBufferChannel),
			"res_buff_chan": len(session.ResPackBufferChannel),
		}
	}

	str, err := json.Marshal(infoMap)
	res = NewGpack(req.Appid, "info", string(str))
	return
}

func GpackAuthHandler(req *Gpack, server *Server) (res *Gpack, err error) {

	if req.Action != "auth" {
		err = errors.New("error")
		return
	}

	key := "app:" + req.Appid + ":baseinfo"
	v, err := server.GameInfoRedisSlave.Get(key)
	if err != nil {
		errorLogger.Println("app base info get error ", err)
		return
	}
	if v == nil {
		errorLogger.Println("app base info get nil ", key)
		err = errors.New("auth: pp base info get nil " + key)
		return
	}

	var gameInfo map[string]interface{}
	if err = json.Unmarshal(v.([]byte), &gameInfo); err != nil {
		errorLogger.Println("game info: data json decoding failed: %s", err)
	}

	data, err := req.GetData()
	if err != nil {
		errorLogger.Println("req data get error ", err)
		return
	}

	if (*data)["app_secret"].(string) == gameInfo["app_secret"].(string) {
		str, err := json.Marshal(map[string]string{"action": "auth",
			"ec":     "0",
			"packid": req.Packid})
		if err != nil {

		}
		res = NewGpack(req.Appid, "ret", string(str))
	} else {
		errorLogger.Println("auth: secret invalid ", key)
		err = errors.New("auth: secret invalid " + key)
	}

	return
}
