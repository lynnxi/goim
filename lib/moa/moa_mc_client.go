package moa

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"game-im/lib/gomemcache"
	"game-im/lib/hash"
	"game-im/lib/stdlog"
	"math/rand"
	"os"
	"time"
)

var (
	moaLogger = stdlog.Log("moa")
)

const (
	KETAMA = 1
	RAND   = 2
)

type MoaMcClient struct {
	hosts       []string
	hash        int
	hashLocator *hash.KetamaHashLocator
	serviceUri  string
	clientCache map[string]*memcache.Client
	time        int64
}

func NewMoaMcClient(hosts []string, serviceUri string, hash int) (client *MoaMcClient) {
	client = &MoaMcClient{hosts: hosts,
		serviceUri:  serviceUri,
		clientCache: make(map[string]*memcache.Client),
		hash:        hash,
		time:        time.Now().Unix()}
	outLogMoa, _ := os.OpenFile("/home/deploy/log/server_moa.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	moaLogger.SetOutput(outLogMoa)

	return
}

func GenerateUid() string {
	b := make([]byte, 5)
	crand.Read(b)
	en := base64.StdEncoding // or URLEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)

	return string(d)
}

func (client *MoaMcClient) IsExpired() bool {
	var t int64 = 5
	now := time.Now().Unix()
	return now-t > client.time
}

func (client *MoaMcClient) Invoke(method string, args []interface{}, hashKey string) (result interface{}, err error) {
	// 构建将要发送到moa的指令
	cmd := map[string]interface{}{
		"id":     GenerateUid(),
		"action": client.serviceUri,
		"params": map[string]interface{}{
			"m":    method,
			"args": args,
		},
	}

	_cmdjson, _ := json.Marshal(cmd)
	cmdjson := string(_cmdjson)
	var respjson []byte
	if len(cmdjson) < 240 {
		moaLogger.Println("invoke short cmd : " + cmdjson)
		respjson, err = client.invokeShortCmd(cmdjson, hashKey)
	} else {
		moaLogger.Println("invoke long cmd : " + cmdjson)
		respjson, err = client.invokeLongCmd(cmdjson, cmd["id"].(string), hashKey)
	}

	if err != nil {
		return
	}
	//moaLogger.Println("moa resp : " + string(respjson))

	var rtn map[string]interface{}
	if err = json.Unmarshal(respjson, &rtn); err != nil {

	}
	if int(rtn["ec"].(float64)) != 0 {
		err = errors.New("moa get error : " + string(respjson))
	}
	result = rtn["result"]

	return
}

func (client *MoaMcClient) invokeShortCmd(cmdjson string, hashKey string) (respjson []byte, err error) {
	item, err := client.getMc(hashKey).Get(cmdjson)
	if err != nil {
		return
	}
	respjson = item.Value

	return
}

func (client *MoaMcClient) invokeLongCmd(cmdjson string, cmdid string, hashKey string) (respjson []byte, err error) {
	mc := client.getMc(hashKey)
	mc.Set(&memcache.Item{Key: "_buf_cmd_", Value: []byte(cmdjson)})
	// 再用get来获取执行结果
	cmdGet, _ := json.Marshal(map[string]interface{}{
		"id":     GenerateUid(),
		"action": "/service/bufcmd",
		"params": map[string]interface{}{
			"m":    "execute",
			"args": [1]string{cmdid},
		},
	})
	item, err := mc.Get(string(cmdGet))
	if err != nil {
		return
	}
	respjson = item.Value

	return
}

func (client *MoaMcClient) getMc(hashKey string) (mc *memcache.Client) {
	var x int
	var host string
	var n = len(client.hosts)

	if n == 0 {

	}

	if n > 1 {
		if client.hash == KETAMA {
			if client.hashLocator == nil {
				client.hashLocator = hash.NewKetamaHashLocator(client.hosts)
			}
			host = client.hashLocator.GetNodeByKey([]byte(hashKey))
		} else {
			x = rand.Intn(n - 1)
			host = client.hosts[x]
		}
	} else {
		host = client.hosts[0]
	}

	moaLogger.Println("use moa host : " + host + " hash key : " + hashKey)
	mc, ok := client.clientCache[host]
	if !ok {
		mc = memcache.New(host)
		client.clientCache[host] = mc
	}

	return
}
