package mio

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"time"
)

type Msg struct {
	To   string //发给谁
	Id   string //消息id
	Text string //消息内容
	Fr   string //发送者
	Type string //消息类型 1-文字
	T    string //时间戳
	V    int    //消息版本号 (lv)
	Ext  string //扩展字段
}

func NewMsg(fr string, to string, text string, typ string, t string, v int, ext string) (msg *Msg) {
	msg = &Msg{
		Id:   GenerateUid(),
		Fr:   fr,
		To:   to,
		Text: text,
		Type: typ,
		T:    t,
		V:    v,
		Ext:  ext,
	}

	return
}

func GetMsgByBytes(bytes []byte) (msg *Msg) {
	if err := json.Unmarshal(bytes, &msg); err != nil {

	}

	return
}

func (msg *Msg) ToBytes() (bytes []byte) {

	bytes, _ = json.Marshal(msg)

	return
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

type Gpack struct {
	Cmd       string
	Sid       string //会话id，标识响应和请求
	Timestamp int64
	Flag      string //标识所属链接
	Body      interface{}
}

func (gpack *Gpack) SetFlag(flag string) {
	gpack.Flag = flag
}

func (gpack *Gpack) GetFlag() string {
	return gpack.Flag
}

func NewGpack(cmd string, sid string, body interface{}) (gpack *Gpack) {
	gpack = &Gpack{Cmd: cmd,
		Body: body}

	if sid != "" {
		gpack.Sid = sid
	} else {
		gpack.Sid = GenerateUid()
	}

	gpack.Timestamp = time.Now().Unix()
	return
}

func NewReqGpack(cmd string, body interface{}) (gpack *Gpack) {
	gpack = NewGpack(cmd, "", body)
	return
}

func NewRetGpack(sid string, body interface{}) (gpack *Gpack) {
	gpack = NewGpack("", sid, body)
	return
}

func GenerateUid() string {
	b := make([]byte, 5)
	rand.Read(b)
	en := base64.StdEncoding // or URLEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)

	return string(d)
}

func GetErrGpack(em string, sid string) *Gpack {
	if sid == "" {
		sid = GenerateUid()
	}
	return NewRetGpack(sid, map[string]string{"ec": "-1", "em": em})
}

func (gpack *Gpack) ToBytes() (bytes []byte) {

	bytes, _ = json.Marshal(gpack)

	return
}

func GpackEncoding(gpack *Gpack) (bytes []byte) {

	bytes, _ = json.Marshal(gpack)

	return
}

func GpackDecoding(bytes []byte) (gpack *Gpack, err error) {

	if err = json.Unmarshal(bytes, &gpack); err != nil {
	}

	return
}
