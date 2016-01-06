package main

import (
	"encoding/json"
	"errors"
	"game-im/config"
	"game-im/lib/mcrypt"
	. "game-im/lib/mio"
	"strconv"
)

/**
 *
 *  消息
 */
func GpackUgmsgHandler(req *Gpack) (res *Gpack, err error) {
	var body map[string]interface{} = req.Body.(map[string]interface{})
	var msg = body["msgs"].([]interface{})[0].(map[string]interface{})
	var unid = msg["To"].(string)
	var text = msg["Text"].(string)
	var appid = req.GetFlag()

	var userid string
	if appid == "ex_mmzb_wFaBpbG" || appid == "ex_mmzbtest_4QGQJKy" {
		userid, err = mcrypt.DecryptV1([]byte(msg["Fr"].(string)), config.SALT)
	} else {
		userid, err = mcrypt.DecryptV2([]byte(msg["Fr"].(string)), []byte(strconv.Itoa(config.SALT)+appid))
	}

	if err != nil || userid == "" {
		err = errors.New("userid decrypt error")
		errorLogger.Println("userid decrypt error : ", err, userid)
		return
	}

	gid, err := Moa.GetGroupByUnionid(appid, unid)

	if err != nil {
		errorLogger.Println("unionid 2 gid get error : ", err)
		return
	}
	if gid == "" {
		str := "unionid to gid get nil : " + req.Sid
		err = errors.New(str)
		noticeLogger.Println(str)
		return
	}

	key := "app:" + appid + ":baseinfo"
	v, err := GameInfoRedisSlave.Get(key)
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

	_msg := map[string]interface{}{
		"from": userid,
		"to":   gid,
		//"to":       "11182817",
		"actions":  "[\"[游戏聊天|goto_app|" + gameInfo["back_url"].(string) + "]\"]",
		"pushText": text,
		"text":     "[" + text + "|n=|lt=0|s=5x5|st=I|apk=A]",
	}

	Moa.PushMsgV3(_msg)

	res = NewRetGpack(req.Sid, map[string]string{"ec": "0"})
	return
}
