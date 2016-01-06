package main

// import (
// 	//"encoding/json"
// 	"errors"
// 	//"game-im/config"
// 	//"game-im/lib/mcrypt"
// 	. "game-im/lib/mio"
// 	//"game-im/lib/stdlog"
// 	// "hash/crc32"
// 	"strconv"
// 	//"encoding/binary"
// )

// func GpackMsgSynHandler(req *Gpack) (res *Gpack, err error) {
// 	lv := int(req.Body.(map[string]interface{})["lv"].(float64))

// 	msgs, err := getMsgs(req.Gameid, lv)
// 	if err != nil {
// 		return
// 	}
// 	mv := getMaxVersion(req.Gameid)
// 	body := map[string]interface{}{
// 		"msgs": msgs,
// 		"lv":   mv,
// 	}

// 	res = NewGpack(req.Gameid, req.Appid, "msg_ack", req.Sid, body)

// 	return
// }

// func GpackMsgFinHandler(req *Gpack) (res *Gpack, err error) {
// 	gameid := req.Gameid

// 	lv := int(req.Body.(map[string]interface{})["lv"].(float64))
// 	mv := getMaxVersion(gameid)
// 	if mv == lv { // 这里不是原子的 有消息延迟的问题
// 		key_psh := "psh:" + gameid + ":flag"
// 		key_msg := "msg:" + gameid + ":queue"
// 		sessionRedis.Del(key_psh)
// 		msgRedis.ZremRangeByScore(key_msg, 0, lv)
// 	}
// 	res = NewRetGpack(req.Gameid, req.Appid, req.Sid, map[string]string{"ec": "0"})

// 	return
// }

// func GpackMsgHandler(req *Gpack) (res *Gpack, err error) {
// 	to := req.Body.(map[string]interface{})["To"].(string)
// 	text := req.Body.(map[string]interface{})["Text"].(string)
// 	ext := req.Body.(map[string]interface{})["Ext"].(string)
// 	//Type := req.Body.(map[string]interface{})["Type"].(string)

// 	version := incrMsgVersion(to)
// 	msg := NewMsg(req.Gameid, to, text, 1, "", version, ext)
// 	setMsg(msg)
// 	sendPsh(to, req.Appid)

// 	res = NewRetGpack(req.Gameid, req.Appid, req.Sid, map[string]string{"ec": "0"})

// 	return
// }

// func getMaxVersion(gameid string) (lv int) {
// 	key := "msg:" + gameid + ":version"
// 	v, err := msgRedis.Get(key)

// 	if err != nil {
// 		errorLogger.Println("msg version get err ", err)
// 		return
// 	}
// 	if v == nil {
// 		errorLogger.Println("msg version get nil ", key)
// 		err = errors.New("msg version get nil " + key)
// 		return
// 	}

// 	lv, _ = strconv.Atoi(string(v.([]byte)))
// 	errorLogger.Println("msg version ", lv)

// 	return
// }

// func getMsgs(gameid string, lv int) (msgs []*Msg, err error) {
// 	key := "msg:" + gameid + ":queue"
// 	v, err := msgRedis.ZrangeByScore(key, lv, -1, true)

// 	for i, b := range v.([]interface{}) {
// 		if i%2 == 1 {
// 		} else {
// 			msg := GetMsgByBytes(b.([]byte))
// 			msgs = append(msgs, msg)

// 			errorLogger.Println("get msg ", string(b.([]byte)))
// 		}
// 	}
// 	return
// }
// func sendPsh(gameid string, appid string) {
// 	key := "psh:" + gameid + ":flag"
// 	v, _ := sessionRedis.Get(key)
// 	if v == nil {
// 		gpack := NewReqGpack(gameid, appid, "msg_psh", map[string]string{})
// 		logic.OutputPackBufferChannel <- gpack
// 		sessionRedis.SetEx(key, "1", 3)
// 	}

// }

// func setMsg(msg *Msg) (err error) {
// 	key := "msg:" + msg.To + ":queue"
// 	err = msgRedis.Zadd(key, string(msg.ToBytes()), msg.V)

// 	if err != nil {
// 		errorLogger.Println("msg version incr err ", err)
// 		return
// 	}
// 	return
// }

// func incrMsgVersion(gameid string) (version int) {
// 	key := "msg:" + gameid + ":version"
// 	v, err := msgRedis.Incr(key)

// 	if err != nil {
// 		errorLogger.Println("msg version incr err ", err)
// 		return
// 	}
// 	if v == nil {
// 		errorLogger.Println("msg version incr nil ", key)
// 		err = errors.New("msg version incr nil " + key)
// 		return
// 	}

// 	version = int(v.(int64))

// 	return
// }
