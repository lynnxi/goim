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

// func GpackGmsgHandler(req *Gpack) (res *Gpack, err error) {
// 	from := req.Gameid
// 	to := req.Body.(map[string]interface{})["To"].(string)
// 	text := req.Body.(map[string]interface{})["Text"].(string)
// 	ext := req.Body.(map[string]interface{})["Ext"].(string)
// 	//Type := req.Body.(map[string]interface{})["Type"].(string)

// 	version := incrGmsgVersion(to)
// 	msg := NewMsg(req.Gameid, to, text, 1, "", version, ext)
// 	setGmsg(msg)
// 	sendGpsh(from, to, req.Appid, version)

// 	res = NewRetGpack(req.Gameid, req.Appid, req.Sid, map[string]string{"ec": "0"})

// 	return
// }

// func GpackGmsgSynHandler(req *Gpack) (res *Gpack, err error) {
// 	lv := int(req.Body.(map[string]interface{})["lv"].(float64))
// 	gid := req.Body.(map[string]interface{})["gid"].(string)

// 	msgs, err := getGmsgs(gid, lv+1)
// 	if err != nil {
// 		return
// 	}
// 	if msgs == nil {
// 		msgs = []*Msg{}
// 	}
// 	mv := getGmaxVersion(gid)

// 	body := map[string]interface{}{
// 		"msgs": msgs,
// 		"lv":   mv,
// 		"gid":  gid,
// 	}

// 	res = NewGpack(req.Gameid, req.Appid, "gmsg_ack", req.Sid, body)

// 	return
// }

// func GpackGmsgFinHandler(req *Gpack) (res *Gpack, err error) {
// 	/**
// 	* @todo 判断版本号继续发psh
// 	**/
// 	//lv := int(req.Body.(map[string]interface{})["lv"].(float64))
// 	gid := req.Body.(map[string]interface{})["gid"].(string)
// 	//mv := getGmaxVersion(gid)
// 	//if mv == lv { // 这里不是原子的 有消息延迟的问题
// 	key := "gpsh:" + gid + ":flag:" + req.Gameid
// 	sessionRedis.Del(key)
// 	//}
// 	res = NewRetGpack(req.Gameid, req.Appid, req.Sid, map[string]string{"ec": "0"})

// 	return
// }

// func getGmsgs(gid string, lv int) (msgs []*Msg, err error) {
// 	key := "gmsg:" + gid + ":queue"
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

// func sendGpsh(from string, gid string, appid string, lv int) {
// 	gameids := getGameidsByGid(gid, appid)
// 	for _, gameid := range gameids {
// 		//if gameid != from {
// 		key := "gpsh:" + gid + ":flag:" + gameid
// 		v, _ := sessionRedis.Get(key)
// 		if v == nil {
// 			gpack := NewReqGpack(gameid, appid, "gmsg_psh", map[string]interface{}{"gid": gid, "lv": lv})
// 			logic.OutputPackBufferChannel <- gpack
// 			sessionRedis.SetEx(key, "1", 3)
// 		}

// 		//}
// 	}
// }

// func setGmsg(msg *Msg) (err error) {
// 	key := "gmsg:" + msg.To + ":queue"
// 	err = msgRedis.Zadd(key, string(msg.ToBytes()), msg.V)

// 	if err != nil {
// 		errorLogger.Println("msg version incr err ", err)
// 		return
// 	}
// 	return
// }

// func getGmaxVersion(gid string) (lv int) {
// 	key := "gmsg:" + gid + ":version"
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

// 	return
// }

// func incrGmsgVersion(gid string) (version int) {
// 	key := "gmsg:" + gid + ":version"
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

// func getGameidsByGid(gid string, appid string) (gameids []string) {
// 	if gid == "1" {
// 		key := "appid:" + appid + ":session"
// 		v, _ := sessionRedis.HgetAll(key)
// 		for i, b := range v.([]interface{}) {
// 			if i%2 == 1 {
// 			} else {
// 				gameid := string(b.([]byte))
// 				gameids = append(gameids, gameid)

// 			}
// 		}
// 	}

// 	return
// }
