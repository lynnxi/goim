package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	//"game-im/lib/stdlog"
)

type Imj struct {
	To     string
	Id     string
	Text   string
	Fr     string
	Action string
	Type   string
	T      string
}

func (imj *Imj) ToString() (str string) {

	str = imj.To + " " + imj.Id + " " + imj.Text + " " + imj.Fr + imj.Action + " " + imj.Type + " " + imj.T

	return
}

type Gpack struct {
	Action string
	Packid string
	Appid  string
	Data   string
}

func NewGpack(appid string, action string, data string) (gpack *Gpack) {
	gpack = &Gpack{Action: action,
		Appid: appid,
		Data:  data}

	gpack.Packid = GenerateUid()
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

var ErrGpack *Gpack

func GetErrGpack() *Gpack {
	if ErrGpack == nil {
		ErrGpack = &Gpack{Action: "ret",
			Appid: "",
			Data:  "{\"ec\":\"-1\"}",
		}

		ErrGpack.Packid = GenerateUid()
	}

	return ErrGpack
}

func (gpack *Gpack) ToString() (str string) {

	str = gpack.Action + " " + gpack.Packid + " " + gpack.Appid + " " + gpack.Data

	return
}

func (gpack *Gpack) GetData() (data *map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(gpack.Data), &data); err != nil {
		errorLogger.Println("gpack: data json decoding failed: %s", err)
	}

	return
}
