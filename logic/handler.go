package main

import (
	//"encoding/json"
	//"errors"
	//"game-im/config"
	//"game-im/lib/mcrypt"
	. "game-im/lib/mio"
	//"game-im/lib/stdlog"
	// "hash/crc32"
	//"strconv"
	//"encoding/binary"
)

var Action2Handler = map[string]func(req *Gpack) (res *Gpack, err error){
	// "msg":      GpackMsgHandler,
	// "msg_syn":  GpackMsgSynHandler,
	// "msg_fin":  GpackMsgFinHandler,
	// "gmsg":     GpackGmsgHandler,
	// "gmsg_syn": GpackGmsgSynHandler,
	// "gmsg_fin": GpackMsgFinHandler,
	// "loc":      GpackLocHandler,
	"ugmsg": GpackUgmsgHandler,
}
