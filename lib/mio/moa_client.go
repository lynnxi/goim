package mio

import (
	"encoding/json"
	"errors"
	"game-im/config"
	"game-im/lib/moa"
	"hash/crc32"
	"strconv"
	"strings"
)

type MoaClient struct {
	clientCache map[string]*moa.MoaRedisClient
}

func NewMoaClient() (client *MoaClient) {
	client = &MoaClient{clientCache: make(map[string]*moa.MoaRedisClient)}
	return
}

func (client *MoaClient) LookUp(serviceUri string, protocol string) ([]string, error) {
	moaMcClient := moa.NewMoaMcClient([]string{config.MOA_LOOKUP_HOST}, config.MOA_LOOKUP_URI, moa.RAND)
	res, err := moaMcClient.Invoke("getService", []interface{}{serviceUri, protocol}, "nil")
	if err != nil {
		return nil, err
	}
	_hosts := res.(map[string]interface{})["hosts"].([]interface{})
	var len = len(_hosts)
	var hosts []string = make([]string, len)
	for i, host := range _hosts { //找不到好的类型转换方法，临时解决
		arr := strings.Split(host.(string), "?")
		hosts[i] = arr[0]
	}
	return hosts, err
}



func (client *MoaClient) GetGroupByUnionid(appid string, unionid string) (gid string, err error) {
	moaRedisClient, ok := client.clientCache[config.MOA_GROUP_PROFILE]

	if !ok || moaRedisClient.IsExpired() {
		hosts, err1 := client.LookUp(config.MOA_GROUP_PROFILE, "redis")

		if err1 != nil {
			err = err1
			return
		}

		moaRedisClient = moa.NewMoaRedisClient(hosts, config.MOA_GROUP_PROFILE, moa.KETAMA)
		client.clientCache[config.MOA_GROUP_PROFILE] = moaRedisClient
	}

	args := []interface{}{appid, unionid}

	res, err := moaRedisClient.Invoke("getGroupByUnionid", args, unionid)
	if err != nil {
		return
	}
	ec := res.(map[interface{}]interface{})["ec"].(int64)
	em := res.(map[interface{}]interface{})["em"].(string)
	if ec != 0 {
		err = errors.New("moa return false ec " + strconv.Itoa(int(ec)) + " em " + em)
	}
	data := res.(map[interface{}]interface{})["data"].(map[interface{}]interface{})
	gid = data["gid"].(string)
	return
}

func (client *MoaClient) GetGroupProfile(gid string, userid string) (profile map[string]interface{}, err error) {
	moaRedisClient, ok := client.clientCache[config.MOA_GROUP_PROFILE]

	if !ok || moaRedisClient.IsExpired() {
		hosts, err1 := client.LookUp(config.MOA_GROUP_PROFILE, "redis")

		if err1 != nil {
			err = err1
			return
		}

		moaRedisClient = moa.NewMoaRedisClient(hosts, config.MOA_GROUP_PROFILE, moa.KETAMA)
		client.clientCache[config.MOA_GROUP_PROFILE] = moaRedisClient
	}

	args := []interface{}{gid, userid}

	res, err := moaRedisClient.Invoke("groupProfile", args, gid)
	if err != nil {
		return
	}
	ec := res.(map[interface{}]interface{})["ec"].(int64)
	em := res.(map[interface{}]interface{})["em"].(string)
	if ec != 0 {
		err = errors.New("moa return false ec " + strconv.Itoa(int(ec)) + " em " + em)
	}
	data := res.(map[interface{}]interface{})["data"].(map[interface{}]interface{})
	profile = map[string]interface{}{
		"name":       data["name"],
		"gid":        data["gid"],
		"game_union": data["game_union"],
	}
	return
}

var (
	CQueueRediss = [...]*RedisClient{
		NewRedisClient(config.REDIS_CQ_1),
		NewRedisClient(config.REDIS_CQ_2),
		NewRedisClient(config.REDIS_CQ_5),
		NewRedisClient(config.REDIS_CQ_6)}
)

func (client *MoaClient) PushMsgV3(msg map[string]interface{}) {
	i := int(crc32.ChecksumIEEE([]byte(msg["to"].(string))))
	i = i % len(CQueueRediss)
	_qredis := CQueueRediss[i]

	command := map[string]interface{}{
		"id":     "",
		"action": "proxy_gmsg",
		"params": map[string]interface{}{
			"m":         "proxyMessage",
			"push_text": msg["pushText"],
			"body":      msg["text"],
			"actions":   msg["actions"],
			"from":      msg["from"],
			"to":        msg["to"],
			"type":      "action",
			"sendtoall": true,
		},
	}

	v, err := json.Marshal(command)
	if err != nil {

	}
	r := strings.NewReplacer("\\\\ue409", "\\ue409", "\\\\n", "\\n")
	s := r.Replace(string(v))
	_qredis.Rpush("_router_v2_command_", s)
}
